package main

func VMHandleECALL(vm VMState) VMState {
	vm.status = VM_STATUS_RUNNING
	return vm
}
