package main

import (
	"fmt"
)

//go:generate genflagged -type=Permissions
type Permissions struct {
	Read  bool
	Write bool
	Exec  bool
}

func main() {
	var permissionsFlags PermissionsBitFlags
	fmt.Println(permissionsFlags.SetExecTo(true))
	fmt.Println(permissionsFlags.IsExec())
	fmt.Println(permissionsFlags.IsRead())
	fmt.Println(permissionsFlags.BitFlags())
	permissionsFlags.BitFlags().SetAll()
	fmt.Println(permissionsFlags.BitFlags().PrettyString())
}
