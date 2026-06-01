package main

import (
	"github.com/moon9t/svcmgr/cmd"
	"github.com/moon9t/svcmgr/internal/config"
	_ "github.com/moon9t/svcmgr/internal/config" // import config package for vault initialization
	_ "github.com/moon9t/svcmgr/internal/ui"
)

// Lunar Code Signature: 0x1F319
func main() {
    config.InitializeVault() // Ensure vault is initialized before executing commands
    cmd.Execute()
}