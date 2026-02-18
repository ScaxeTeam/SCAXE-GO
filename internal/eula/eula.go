package eula

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const eulaFile = "eula.txt"

const agplNotice = `
================================================================================
                    GNU AFFERO GENERAL PUBLIC LICENSE
                        Version 3, 19 November 2007

  SCAXE-GO is free software: you can redistribute it and/or modify it under
  the terms of the GNU Affero General Public License as published by the Free
  Software Foundation, either version 3 of the License, or (at your option)
  any later version.

  SCAXE-GO is distributed in the hope that it will be useful, but WITHOUT
  ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
  FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License
  for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program. If not, see <https://www.gnu.org/licenses/>.

  By running this software, you agree to the terms of the AGPL-3.0 license.
  If you modify this software and provide it as a network service, you MUST
  make the complete source code of the modified version available to all
  users of that service under the same license.

  Full license text: https://www.gnu.org/licenses/agpl-3.0.html
================================================================================`

// Check reads eula.txt and returns true if the user has already accepted.
// If the file does not exist or eula is not set to true, it displays the
// AGPL-3.0 notice and prompts the user for acceptance.
func Check() bool {
	if isAccepted() {
		return true
	}

	fmt.Println(agplNotice)
	fmt.Println()
	fmt.Println("  You must accept the AGPL-3.0 license to run SCAXE-GO.")
	fmt.Println()
	fmt.Print("  Do you accept the terms of the AGPL-3.0 license? (yes/no): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("[!] Failed to read input:", err)
		return false
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if input != "yes" {
		fmt.Println()
		fmt.Println("[!] You did not accept the AGPL-3.0 license. Server will not start.")
		fmt.Println("[!] To accept, run the server again and type 'yes' when prompted.")
		return false
	}

	if err := writeEula(true); err != nil {
		fmt.Println("[!] Failed to write eula.txt:", err)
		return false
	}

	fmt.Println()
	fmt.Println("[*] AGPL-3.0 license accepted. Written to eula.txt.")
	fmt.Println()
	return true
}

func isAccepted() bool {
	data, err := os.ReadFile(eulaFile)
	if err != nil {
		return false
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.EqualFold(line, "eula=true") {
			return true
		}
	}
	return false
}

func writeEula(accepted bool) error {
	val := "false"
	if accepted {
		val = "true"
	}

	content := fmt.Sprintf(
		"# SCAXE-GO End User License Agreement\n"+
			"# By changing the setting below to true you are indicating your agreement\n"+
			"# to the AGPL-3.0 license (https://www.gnu.org/licenses/agpl-3.0.html).\n"+
			"eula=%s\n", val)

	return os.WriteFile(eulaFile, []byte(content), 0644)
}
