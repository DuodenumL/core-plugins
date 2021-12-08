package schedule

import (
	"container/heap"
	"math"
	"sort"

	"github.com/projecteru2/core-plugins/volume/types"
)

type volume struct {
	device string
	size   int64
}

func (v *volume) LessThan(v1 *volume) bool {
	if v.size == v1.size {
		return v.device < v1.device
	}
	return v.size < v1.size
}

type volumes []*volume

// DeepCopy .
func (v volumes) DeepCopy() volumes {
	res := volumes{}
	for _, item := range v {
		res = append(res, &volume{device: item.device, size: item.size})
	}
	return res
}

type volumeHeap volumes

// Len .
func (v volumeHeap) Len() int {
	return len(v)
}

// Less .
func (v volumeHeap) Less(i, j int) bool {
	return v[i].LessThan(v[j])
}

// Swap .
func (v volumeHeap) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// Push .
func (v *volumeHeap) Push(x interface{}) {
	*v = append(*v, x.(*volume))
}

// Pop .
func (v *volumeHeap) Pop() interface{} {
	old := *v
	n := len(old)
	x := old[n-1]
	*v = old[:n-1]
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type host struct {
	maxDeployCount int
	usedVolumes    volumes
	unusedVolumes  volumes
}

func newHost(resourceInfo *types.NodeResourceInfo, maxDeployCount int) *host {
	h := &host{
		maxDeployCount: maxDeployCount,
		usedVolumes:    []*volume{},
		unusedVolumes:  []*volume{},
	}

	for device, size := range resourceInfo.Capacity.Volumes {
		used := resourceInfo.Usage.Volumes[device]
		if used == 0 {
			h.unusedVolumes = append(h.unusedVolumes, &volume{device: device, size: size})
		} else {
			h.usedVolumes = append(h.usedVolumes, &volume{device: device, size: size - used})
		}
	}

	sort.SliceStable(h.unusedVolumes, func(i, j int) bool { return h.unusedVolumes[i].LessThan(h.unusedVolumes[j]) })
	sort.SliceStable(h.usedVolumes, func(i, j int) bool { return h.usedVolumes[i].LessThan(h.usedVolumes[j]) })
	return h
}

func (h *host) getMonoPlans(monoRequests types.VolumeBindings) ([]types.VolumePlan, int) {
	if len(monoRequests) == 0 {
		return []types.VolumePlan{}, h.maxDeployCount
	}
	if len(h.unusedVolumes) == 0 {
		return []types.VolumePlan{}, 0
	}

	volumes := h.unusedVolumes.DeepCopy()
	volumePlans := []types.VolumePlan{}
	volumePlan := types.VolumePlan{}

	// h.unusedVolumes and monoRequests have already been sorted
	reqIndex := 0
	for volumeIndex := 0; volumeIndex < len(volumes); volumeIndex++ {
		request := monoRequests[reqIndex]
		volume := volumes[volumeIndex]
		if volume.size < request.SizeInBytes {
			continue
		}
		volumePlan[request] = types.VolumeMap{volume.device: volume.size}
		if reqIndex == len(monoRequests)-1 {
			volumePlans = append(volumePlans, volumePlan)
			volumePlan = types.VolumePlan{}
		}
		reqIndex = (reqIndex + 1) % len(monoRequests)
	}
	return volumePlans, len(volumePlans)
}

func (h *host) getNormalPlan(vHeap *volumeHeap, normalRequests types.VolumeBindings) types.VolumePlan {
	volumePlan := types.VolumePlan{}
	for reqIndex := 0; reqIndex < len(normalRequests); reqIndex++ {
		req := normalRequests[reqIndex]
		volumeToPush := []*volume{}
		allocated := false

		for vHeap.Len() > 0 {
			volume := heap.Pop(vHeap).(*volume)
			if volume.size >= req.SizeInBytes {
				volumePlan[req] = types.VolumeMap{volume.device: req.SizeInBytes}
				allocated = true
				volume.size -= req.SizeInBytes
				if volume.size > 0 {
					volumeToPush = append(volumeToPush, volume)
				}
				break
			}
		}

		for _, volume := range volumeToPush {
			heap.Push(vHeap, volume)
		}

		if !allocated {
			return nil
		}
	}

	return volumePlan
}

func (h *host) getNormalPlans(normalRequests types.VolumeBindings) ([]types.VolumePlan, int) {
	if len(normalRequests) == 0 {
		return []types.VolumePlan{}, h.maxDeployCount
	}

	vh := volumeHeap(h.usedVolumes.DeepCopy())
	vHeap := &vh
	heap.Init(vHeap)

	volumePlans := []types.VolumePlan{}

	for len(volumePlans) <= h.maxDeployCount {
		if volumePlan := h.getNormalPlan(vHeap, normalRequests); volumePlan != nil {
			volumePlans = append(volumePlans, volumePlan)
		} else {
			break
		}
	}
	return volumePlans, len(volumePlans)
}

func (h *host) classifyVolumeBindings(volumeBindings types.VolumeBindings) (normalRequests, monoRequests, unlimitedRequests types.VolumeBindings) {
	for _, binding := range volumeBindings {
		switch {
		case binding.RequireScheduleMonopoly():
			monoRequests = append(monoRequests, binding)
		case binding.RequireScheduleUnlimitedQuota():
			unlimitedRequests = append(unlimitedRequests, binding)
		case binding.RequireSchedule():
			normalRequests = append(normalRequests, binding)
		}
	}

	sort.SliceStable(monoRequests, func(i, j int) bool { return monoRequests[i].SizeInBytes < monoRequests[j].SizeInBytes })
	sort.SliceStable(normalRequests, func(i, j int) bool { return normalRequests[i].SizeInBytes < normalRequests[j].SizeInBytes })

	return normalRequests, monoRequests, unlimitedRequests
}

func (h *host) getUnlimitedPlans(normalPlans, monoPlans []types.VolumePlan, unlimitedRequests types.VolumeBindings) []types.VolumePlan {
	capacity := len(normalPlans)

	volumes := append(h.usedVolumes.DeepCopy(), h.unusedVolumes.DeepCopy()...)
	volumeMap := map[string]*volume{}
	for _, volume := range volumes {
		volumeMap[volume.device] = volume
	}

	// apply changes
	for _, plan := range normalPlans {
		for _, vm := range plan {
			for device, size := range vm {
				volumeMap[device].size -= size
			}
		}
	}
	for _, plan := range monoPlans {
		for _, vm := range plan {
			for device, size := range vm {
				volumeMap[device].size -= size
			}
		}
	}

	// select the volume with the largest size
	maxVolume := volumes[0]
	for i := range volumes {
		if volumes[i].size > maxVolume.size {
			maxVolume = volumes[i]
		}
	}

	plans := []types.VolumePlan{}
	for i := 0; i < capacity; i++ {
		plan := types.VolumePlan{}
		for _, req := range unlimitedRequests {
			plan[req] = types.VolumeMap{maxVolume.device: req.SizeInBytes}
		}
		plans = append(plans, plan)
	}

	return plans
}

func (h *host) getVolumePlans(volumeBindings types.VolumeBindings) []types.VolumePlan {
	if len(h.unusedVolumes) == 0 && len(h.usedVolumes) == 0 {
		return nil
	}

	normalRequests, monoRequests, unlimitedRequests := h.classifyVolumeBindings(volumeBindings)

	minNormalRequestSize := int64(math.MaxInt)
	for _, normalRequest := range normalRequests {
		if normalRequest.SizeInBytes < minNormalRequestSize {
			minNormalRequestSize = normalRequest.SizeInBytes
		}
	}

	// get baseline
	normalPlans, normalCapacity := h.getNormalPlans(normalRequests)
	monoPlans, monoCapacity := h.getMonoPlans(monoRequests)
	bestCapacity := min(monoCapacity, normalCapacity)
	bestVolumePlans := [2][]types.VolumePlan{normalPlans[:min(bestCapacity, len(normalPlans))], monoPlans[:min(bestCapacity, len(monoPlans))]}

	for len(monoPlans) > len(normalPlans) && len(h.unusedVolumes) >= len(monoRequests) {
		// convert an unused volume to used volume
		p := sort.Search(len(h.unusedVolumes), func(i int) bool { return h.unusedVolumes[i].size >= minNormalRequestSize })
		// if no volume meets the requirement, just stop
		if p == len(h.unusedVolumes) {
			break
		}
		v := h.unusedVolumes[p]
		h.unusedVolumes = append(h.unusedVolumes[:p], h.unusedVolumes[p+1:]...)
		h.usedVolumes = append(h.usedVolumes, v)

		normalPlans, normalCapacity = h.getNormalPlans(normalRequests)
		monoPlans, monoCapacity = h.getMonoPlans(monoRequests)
		capacity := min(monoCapacity, normalCapacity)
		if capacity > bestCapacity {
			bestCapacity = capacity
			bestVolumePlans = [2][]types.VolumePlan{normalPlans[:capacity], monoPlans[:capacity]}
		}
	}

	normalPlans, monoPlans = bestVolumePlans[0], bestVolumePlans[1]
	unlimitedPlans := h.getUnlimitedPlans(normalPlans, monoPlans, unlimitedRequests)

	resVolumePlans := []types.VolumePlan{}
	merge := func(plan types.VolumePlan, plans []types.VolumePlan, index int) {
		if index < len(plans) {
			plan.Merge(plans[index])
		}
	}

	for i := 0; i < bestCapacity; i++ {
		plan := types.VolumePlan{}
		merge(plan, normalPlans, i)
		merge(plan, monoPlans, i)
		merge(plan, unlimitedPlans, i)
		resVolumePlans = append(resVolumePlans, plan)
	}
	return resVolumePlans
}

// GetVolumePlans .
func GetVolumePlans(resourceInfo *types.NodeResourceInfo, volumeRequest types.VolumeBindings, existing types.VolumeMap, maxDeployCount int) []types.VolumePlan {
	if existing != nil {
		// todo: affinity
	}

	h := newHost(resourceInfo, maxDeployCount)
	return h.getVolumePlans(volumeRequest)
}
