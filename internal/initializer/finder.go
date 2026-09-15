package initializer

import (
	"fmt"
	"os"
	"path/filepath"
)

const workflowXML = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AMApplicationBuild</key>
	<string>521.1</string>
	<key>AMApplicationVersion</key>
	<string>2.14</string>
	<key>AMDocumentVersion</key>
	<string>2</string>
	<key>actions</key>
	<array>
		<dict>
			<key>action</key>
			<dict>
				<key>AMActionVersion</key>
				<string>2.0.3</string>
				<key>AMApplication</key>
				<string>Finder</string>
				<key>AMParameterProperties</key>
				<dict>
					<key>COMMAND_STRING</key>
					<dict/>
					<key>CheckedForUserDefaultShell</key>
					<dict/>
					<key>inputMethod</key>
					<dict/>
					<key>shell</key>
					<dict/>
					<key>source</key>
					<dict/>
				</dict>
				<key>AMProvides</key>
				<dict>
					<key>Container</key>
					<string>List</string>
					<string>Types</string>
					<string>Text</string>
				</dict>
				<key>ActionBundlePath</key>
				<string>/System/Library/Automator/Run Shell Script.action</string>
				<key>ActionName</key>
				<string>Run Shell Script</string>
				<key>ActionParameters</key>
				<dict>
					<key>COMMAND_STRING</key>
					<string>for f in "$@"
do
    cd "$f" || exit
    /opt/homebrew/bin/context --clipboard
done</string>
					<key>CheckedForUserDefaultShell</key>
					<true/>
					<key>inputMethod</key>
					<integer>1</integer>
					<key>shell</key>
					<string>/bin/bash</string>
					<key>source</key>
					<string></string>
				<key>BundleIdentifier</key>
				<string>com.apple.Automator.RunShellScript</string>
				<key>CFBundleVersion</key>
				<string>2.0.3</string>
				<key>Canosectic/ig/success</key>
				<string></string>
			</dict>
		</dict>
	</array>
	<key>connectors</key>
	<dict/>
	<key>workflowMetaData</key>
	<dict>
		<key>applicationBundleID</key>
		<string>com.apple.finder</string>
		<key>applicationBundleIDsByPath</key>
		<dict>
			<key>string>/System/Library/Finder.app</key>
		</dict>
		<key>applicationName</key>
		<string>Finder</string>
		<key>inputTypeIdentifier</key>
		<string>com.apple.Automator.fileSystemObject.folder</string>
		<key>outputTypeIdentifier</key>
		<string>com.apple.Automator.nothing</string>
		<key>presentationMode</key>
		<string>15</string>
		<key>processesInputIcon</key>
		<string>Context</string>
	</dict>
</dict>
</plist>
`

func InstallFinderIntegration() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	workflowPath := filepath.Join(home, "Library", "Services", "Generate AI Context.workflow", "Contents")
	if err := os.MkdirAll(workflowPath, 0755); err != nil {
		return err
	}

	plistFile := filepath.Join(workflowPath, "document.wflow")
	if err := os.WriteFile(plistFile, []byte(workflowXML), 0644); err != nil {
		return err
	}

	fmt.Println("✓ Finder context menu integration installed successfully!")
	fmt.Println("  You can now right-click any folder in Finder -> Quick Actions -> Generate AI Context")
	return nil
}
