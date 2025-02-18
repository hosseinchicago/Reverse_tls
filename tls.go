//Author:github.com/Azumi67
//This script is for educational use and for my own learning, but I'd be happy if you find it useful too.
//This script simplifies the configuration of WS + WSS Reverse tunnel.
//You can send me feedback so I can use it to learn more.
//This script comes without any warranty
//Thank you.
package main

import (
    "bytes"
    "time"
	"strconv"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/fatih/color"
	"log"
	"github.com/AlecAivazis/survey/v2"
	"net"
	"io/ioutil"
)
func getIPv4() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		name := iface.Name
		if strings.HasPrefix(name, "eth") || strings.HasPrefix(name, "en") {
			addresses, err := iface.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addresses {
				ipnet, ok := addr.(*net.IPNet)
				if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}

	return ""
}
func displayProgress(total, current int) {
	width := 40
	percentage := current * 100 / total
	completed := width * current / total
	remaining := width - completed

	fmt.Printf("\r[%s>%s] %d%%", strings.Repeat("=", completed), strings.Repeat(" ", remaining), percentage)
}

func displayError(message string) {
	fmt.Printf("\u2718 Error: %s\n", message)
}

func displayNotification(message string) {
	fmt.Printf("\033[93m%s\033[0m\n", message)
}

func displayCheckmark(message string) {
	fmt.Printf("\033[92m\u2714 \033[0m%s\n", message)
}

func displayLoading() {
    frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
    delay := 100 * time.Millisecond
    duration := 5 * time.Second

    endTime := time.Now().Add(duration)

    for time.Now().Before(endTime) {
        for _, frame := range frames {
            fmt.Printf("\r[%s] Loading... ", frame)
            time.Sleep(delay)
        }
    }
    fmt.Println()
}
func displayLogo2() error {
	cmd := exec.Command("bash", "-c", "/etc/./logo.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
func displayLogo() {
	logo := `
   ______    _______    __      _______          __       _____  ___  
  /    " \  |   __ "\  |" \    /"      \        /""\      (\"   \|" \ 
 // ____  \ (. |__) :) ||  |  |:        |      /    \     |.\\   \   |
/  /    ) :)|:  ____/  |:  |  |_____/   )     /' /\  \    |: \.   \\ |
(: (____/ // (|  /     |.  |   //       /    //  __'  \   |.  \    \ |
\        // |__/ \     /\  |\  |:  __   \   /   /  \\  \  |    \    \|
 \"_____ / (_______)  (__\_|_) |__|  \___) (___/    \___) \___|\____\)
`
	
    cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
    blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()      
    yellow := color.New(color.FgHiYellow, color.Bold).SprintFunc()   
    red := color.New(color.FgHiRed, color.Bold).SprintFunc()        


	

	    logo = cyan("  ______   ") + blue(" _______  ") + green("  __    ") + yellow("   ______   ") + red("     __      ") + cyan("  _____  ___  \n") +
		cyan(" /     \" \\  ") + blue("|   __ \" ") + green(" |\" \\  ") + yellow("   /\"     \\   ") + red("   /\"\"\\     ") + cyan(" (\\\"   \\|\"  \\ \n") +
		cyan("//  ____  \\ ")  + blue("(. |__) :)") + green("||  |  ") + yellow(" |:       |  ") + red("  /    \\   ") + cyan("  |.\\\\   \\   |\n") +
		cyan("/  /    ) :)") + blue("|:  ____/ ") + green("|:  |  ") + yellow(" |_____/  )  ") + red(" /' /\\  \\   ") + cyan(" |: \\.   \\\\ |\n") +
		cyan("(: (____/ / ") + blue("(|  /     ") + green("|.  | ") + yellow("  //      /  ") + red("//   __'  \\  ") + cyan(" |.  \\    \\ |\n") +
		cyan("\\        / ") + blue("/|__/ \\   ") + green(" /\\  |\\ ") + yellow(" |:  __  \\ ") + red(" /   /  \\\\   ") + cyan ("  |    \\    \\|\n") +
		cyan(" \"_____ / ") + blue("(_______)") + green("  (__\\_|_)") + yellow(" (__) \\___)") + red("(___/    \\___)") + cyan(" \\___|\\____\\)\n")


	fmt.Println(logo)
}
func main() {
	if os.Geteuid() != 0 {
		fmt.Println("\033[91mThis script must be run as root. Please use sudo -i.\033[0m")
		os.Exit(1)
	}

	mainMenu()
}
func readInput() {
	fmt.Print("Press Enter to continue..")
	fmt.Scanln()
	mainMenu()
}
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func mainMenu() {
	for {
		err := displayLogo2()
		if err != nil {
			log.Fatalf("failed to display logo: %v", err)
		}
		displayLogo()
		border := "\033[93m+" + strings.Repeat("=", 70) + "+\033[0m"
		content := "\033[93m║            ▌║█║▌│║▌│║▌║▌█║ \033[92mMain Menu\033[93m  ▌│║▌║▌│║║▌█║▌                  ║"
		footer := " \033[92m            Join Opiran Telegram \033[34m@https://t.me/OPIranClub\033[0m "

		borderLength := len(border) - 2
		centeredContent := fmt.Sprintf("%[1]*s", -borderLength, content)

		fmt.Println(border)
		fmt.Println(centeredContent)
		fmt.Println(border)

		fmt.Println(border)
		fmt.Println(footer)
		fmt.Println(border)
		prompt := &survey.Select{
			Message: "Enter your choice Please:",
			Options: []string{"0. \033[91mSTATUS Menu\033[0m", "1. \033[96mWS \033[92mTCP \033[0m", "2. \033[93mWS \033[92mUDP\033[0m", "3. \033[96mWSS + TLS \033[92mTCP\033[0m", "4. \033[93mWSS + TLS \033[92mUDP\033[0m", "5. \033[92mStop | Restart Service\033[0m", "6. \033[91mUninstall\033[0m", "q. Exit"},
		
		}
		fmt.Println("\033[93m╰─────────────────────────────────────────────────────────────────────╯\033[0m")

		var choice string
		err = survey.AskOne(prompt, &choice)
		if err != nil {
			log.Fatalf("\033[91muser input is wrong:\033[0m %v", err)
		}
		switch choice {
		case "0. \033[91mSTATUS Menu\033[0m":
			status()
		case "1. \033[96mWS \033[92mTCP \033[0m":
			wsMenu1()
		case "2. \033[93mWS \033[92mUDP\033[0m":
			wsMenu2()
		case "3. \033[96mWSS + TLS \033[92mTCP\033[0m":
			wsMenu3()
		case "4. \033[93mWSS + TLS \033[92mUDP\033[0m":
			wsMenu4()
		case "5. \033[92mStop | Restart Service\033[0m":
			startMain()
		case "6. \033[91mUninstall\033[0m":
			UniMenu()
		case "q. Exit":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice.")
		}

		
		readInput()
	}
}
func rmv() error {
	file := "/etc/tls.sh"
	if _, err := os.Stat(file); err == nil {
		err := os.Remove(file)
		if err != nil {
			return fmt.Errorf("\033[91mbash file doesn't exists:\033[0m %v", err)
		}
		fmt.Println("\033[91mbash file removed successfully!\033[0m")
	}
	return nil
}
func deleteCron() {
	entriesToDelete := []string{
		"0 */1 * * * /etc/tls.sh",
		"0 */2 * * * /etc/tls.sh",
		"0 */3 * * * /etc/tls.sh",
		"0 */4 * * * /etc/tls.sh",
		"0 */5 * * * /etc/tls.sh",
		"0 */6 * * * /etc/tls.sh",
		"0 */7 * * * /etc/tls.sh",
		"0 */8 * * * /etc/tls.sh",
		"0 */9 * * * /etc/tls.sh",
		"0 */10 * * * /etc/tls.sh",
		"0 */11 * * * /etc/tls.sh",
		"0 */12 * * * /etc/tls.sh",
		"0 */13 * * * /etc/tls.sh",
		"0 */14 * * * /etc/tls.sh",
		"0 */15 * * * /etc/tls.sh",
		"0 */16 * * * /etc/tls.sh",
		"0 */17 * * * /etc/tls.sh",
		"0 */18 * * * /etc/tls.sh",
		"0 */19 * * * /etc/tls.sh",
		"0 */20 * * * /etc/tls.sh",
		"0 */21 * * * /etc/tls.sh",
		"0 */22 * * * /etc/tls.sh",
		"0 */23 * * * /etc/tls.sh",
		"0 */24 * * * /etc/tls.sh",
	}

	existingCrontab, err := exec.Command("crontab", "-l").Output()
	if err != nil {
		displayError("\033[91mNo existing cron found!\033[0m")
		return
	}

	newCrontab := string(existingCrontab)
	for _, entry := range entriesToDelete {
		if strings.Contains(newCrontab, entry) {
			newCrontab = strings.Replace(newCrontab, entry, "", -1)
		}
	}

	if newCrontab != string(existingCrontab) {
		cmd := exec.Command("crontab")
		cmd.Stdin = strings.NewReader(newCrontab)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		displayNotification("\033[92mDeleting Previous Crons..\033[0m")
	} else {
		displayError("\033[91mNothing Found, moving on..!\033[0m")
	}
}
func resKharej() {
	deleteCron()

	if _, err := os.Stat("/etc/tls.sh"); err == nil {
		os.Remove("/etc/tls.sh")
	}

	file, err := os.Create("/etc/tls.sh")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString("#!/bin/bash\n")
	file.WriteString("systemctl daemon-reload\n")
	file.WriteString("systemctl restart rtun-kharej\n")

	cmd := exec.Command("chmod", "+x", "/etc/tls.sh")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\033[93m╭──────────────────────────────────────╮\033[0m")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\033[93mEnter the \033[92mReset timer:\033[0m ")
	hoursStr, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	hoursStr = strings.TrimSpace(hoursStr)
	fmt.Println("\033[93m╰──────────────────────────────────────╯\033[0m")

	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		displayError("\033[91mInvalid input for the reset timer!\033[0m")
		return
	}

	cronEntry := fmt.Sprintf("0 */%d * * * /etc/tls.sh", hours)
	existingCrontab, err := exec.Command("crontab", "-l").Output()
	if err != nil {
		displayError("\033[91mNo cron found!\033[0m")
		return
	}

	newCrontab := fmt.Sprintf("%s\n%s\n", existingCrontab, cronEntry)
	cmd = exec.Command("crontab")
	cmd.Stdin = strings.NewReader(newCrontab)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("\033[91mUpdate cron error: %v. error result: %s\033[0m", err, stderr.String())
	}
}
func resIran() {
	deleteCron()

	if _, err := os.Stat("/etc/tls.sh"); err == nil {
		os.Remove("/etc/tls.sh")
	}

	file, err := os.Create("/etc/tls.sh")
	if err != nil {
		log.Fatalf("\033[91mCreating bash file error!: %v\033[0m", err)
	}
	defer file.Close()

	file.WriteString("#!/bin/bash\n")
	file.WriteString("systemctl daemon-reload\n")
	file.WriteString("systemctl restart rtun-iran\n")

	cmd := exec.Command("chmod", "+x", "/etc/tls.sh")
	if err := cmd.Run(); err != nil {
		log.Fatalf("\033[91mchmod cmd error!: %v\033[0m", err)
	}

	fmt.Println("\033[93m╭──────────────────────────────────────╮\033[0m")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\033[93mEnter the \033[92mReset timer:\033[0m ")
	hoursStr, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("\033[91mError reading input: %v\033[0m", err)
	}
	hoursStr = strings.TrimSpace(hoursStr)
	fmt.Println("\033[93m╰──────────────────────────────────────╯\033[0m")

	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		log.Fatalf("\033[91mInvalid input for the reset timer: %v\033[0m", err)
	}

	cronEntry := fmt.Sprintf("0 */%d * * * /etc/tls.sh", hours)
	existingCrontab, err := exec.Command("crontab", "-l").Output()
	if err != nil {
		log.Fatalf("\033[91mRetrieve cron error!: %v\033[0m", err)
	}

	newCrontab := fmt.Sprintf("%s\n%s\n", existingCrontab, cronEntry)

	cmd = exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(newCrontab)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("\033[91mUpdate cron error: %v. error resultt: %s\033[0m", err, stderr.String())
	}
}
func runCmd(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Couldn't Run Command '%s': %v", command, err)
	}
}
func startMain() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[92m Service \033[93mMenu\033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mRestart\033[0m", "2. \033[93mStop \033[0m", "0. \033[94mBack to the main menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mRestart\033[0m":
		start()
	case "2. \033[93mStop \033[0m":
		stop()
	case "0. \033[94mBack to the previous menu\033[0m":
	    clearScreen()
		mainMenu()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
func start() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[92m Restart \033[93mMenu\033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mWS\033[0m", "2. \033[93mWSS \033[0m", "0. \033[94mBack to the previous menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mWS\033[0m":
		restartws()
	case "2. \033[93mWS \033[0m":
		restartws()
	case "0. \033[94mBack to the previous menu\033[0m":
	    clearScreen()
		startMain()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
func restartws() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	displayNotification("\033[93mRestarting WS Reverse \033[93m..\033[0m")
	fmt.Println("\033[93m╭─────────────────────────────────────────────╮\033[0m")

	cmd = exec.Command("systemctl", "restart", "rtun-kharej")
	cmd.Run()
	time.Sleep(1 * time.Second)

	cmd = exec.Command("systemctl", "restart", "rtun-iran")
	cmd.Run()
	time.Sleep(1 * time.Second)

	fmt.Print("Progress: ")

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	delay := 0.1
	duration := 1.0
	endTime := time.Now().Add(time.Duration(duration) * time.Second)

	for time.Now().Before(endTime) {
		for _, frame := range frames {
			fmt.Printf("\r[%s] Loading...  ", frame)
			time.Sleep(time.Duration(delay * float64(time.Second)))
			fmt.Printf("\r[%s]             ", frame)
			time.Sleep(time.Duration(delay * float64(time.Second)))
		}
	}

	displayCheckmark("\033[92mRestart completed!\033[0m")
}
func stop() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[92m Stop \033[93mMenu\033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mWS\033[0m", "2. \033[93mWSS \033[0m", "0. \033[94mBack to the previous menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mWS\033[0m":
		stopws()
	case "2. \033[93mWSS \033[0m":
		stopws()
	case "0. \033[94mBack to the previous menu\033[0m":
	    clearScreen()
		startMain()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
func stopws() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	displayNotification("\033[93mStopping WS Reverse \033[93m..\033[0m")
	fmt.Println("\033[93m╭─────────────────────────────────────────────╮\033[0m")

	cmd = exec.Command("systemctl", "stop", "rtun-kharej")
	cmd.Run()
	time.Sleep(1 * time.Second)

	cmd = exec.Command("systemctl", "stop", "rtun-iran")
	cmd.Run()
	time.Sleep(1 * time.Second)

	fmt.Print("Progress: ")

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	delay := 0.1
	duration := 1.0
	endTime := time.Now().Add(time.Duration(duration) * time.Second)

	for time.Now().Before(endTime) {
		for _, frame := range frames {
			fmt.Printf("\r[%s] Loading...  ", frame)
			time.Sleep(time.Duration(delay * float64(time.Second)))
			fmt.Printf("\r[%s]             ", frame)
			time.Sleep(time.Duration(delay * float64(time.Second)))
		}
	}

	displayCheckmark("\033[92mService Stopped!\033[0m")
}
func status() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[92m Status \033[93mMenu\033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mWS\033[0m", "2. \033[93mWSS \033[0m", "0. \033[94mBack to the main menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mWS\033[0m":
		wsStatus()
	case "2. \033[93mWSS \033[0m":
		wsStatus()
	case "0. \033[94mBack to the main menu\033[0m":
	    clearScreen()
		mainMenu()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
func wsStatus() {
	services := []string{"rtun-iran", "rtun-kharej"}

	fmt.Println("\033[93m            ╔════════════════════════════════════════════╗\033[0m")
	fmt.Println("\033[93m            ║               \033[92mReverse Status\033[93m               ║\033[0m")
	fmt.Println("\033[93m            ╠════════════════════════════════════════════╣\033[0m")

	for _, service := range services {
		cmd := exec.Command("systemctl", "is-active", "--quiet", service)
		err := cmd.Run()
		if err != nil {
			continue
		}

		status := "\033[92m✓ Active      \033[0m"
		displayName := ""
		switch service {
		case "rtun-iran":
			displayName = "\033[93mIRAN Server   \033[0m"
		case "rtun-kharej":
			displayName = "\033[93mKharej Server \033[0m"
		default:
			displayName = service
		}

		fmt.Printf("           \033[93m ║\033[0m    %s   |    %s\033[93m    ║\033[0m\n", displayName, status)
	}

	fmt.Println("\033[93m            ╚════════════════════════════════════════════╝\033[0m")
}
func UniMenu() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Uninstallation \033[96mMenu\033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mWS\033[0m", "2. \033[93mWSS \033[0m", "0. \033[94mBack to the main menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mWS\033[0m":
		removews()
	case "2. \033[93mWSS \033[0m":
		removews()
	case "0. \033[94mBack to the main menu\033[0m":
	    clearScreen()
		mainMenu()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
func removews() {
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	displayNotification("\033[93mRemoving Config ..\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	deleteCron()
	rmv()


	if _, err := os.Stat("/root/reverse-tunnel"); err == nil {
		os.RemoveAll("/root/reverse-tunnel")
	}

	azumiServices := []string{
		"rtun-iran", "rtun-kharej",
	}

	for _, serviceName := range azumiServices {
		hideCmd("systemctl", "disable", serviceName+".service")
		hideCmd("systemctl", "stop", serviceName+".service")
		hideCmd("rm", "/etc/systemd/system/"+serviceName+".service")
	}

	runCmd("systemctl", "daemon-reload")

	fmt.Print("Progress: ")

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	delay := 100 * time.Millisecond
	duration := 1 * time.Second
	endTime := time.Now().Add(duration)

	for time.Now().Before(endTime) {
		for _, frame := range frames {
			fmt.Printf("\r[%s] Loading...  ", frame)
			time.Sleep(delay)
			fmt.Printf("\r[%s]             ", frame)
			time.Sleep(delay)
		}
	}

	displayCheckmark("\033[92m Uninstallation completed!\033[0m")
}
func hideCmd(cmd string, args ...string) error {
	command := exec.Command(cmd, args...)

	nullDevice, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	command.Stdout = nullDevice
	command.Stderr = nullDevice

	return command.Run()
}
func install1() {
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	displayNotification("Installing Reverse Tunnel")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	hideCmd("apt", "install", "git", "-y")
	runCmd("git", "clone", "https://github.com/snsinfu/reverse-tunnel")
	
	err := os.Chdir("reverse-tunnel")
	if err != nil {
		log.Fatalf("Couldn't Change Dir: %v", err)
	}
	runCmd("go", "build", "-o", "rtun-server", "github.com/snsinfu/reverse-tunnel/server/cmd")
	displayCheckmark(fmt.Sprintf("\033[92mInstallation was Successful\033[0m"))
}

func install2() {
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	displayNotification("Installing Reverse Tunnel")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	hideCmd("apt", "install", "git", "-y")
	runCmd("git", "clone", "https://github.com/snsinfu/reverse-tunnel")
	
	err := os.Chdir("reverse-tunnel")
	if err != nil {
		log.Fatalf("Couldn't Change Dir: %v", err)
	}
	runCmd("go", "build", "-o", "rtun-client", "github.com/snsinfu/reverse-tunnel/agent/cmd")
	displayCheckmark(fmt.Sprintf("\033[92mInstallation was Successful\033[0m"))
}
func wsMenu1() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Reverse \033[92mWS \033[96mTCP\033[0m ")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mIPV4\033[0m", "2. \033[93mIPV6\033[0m", "0. \033[94mBack to the main menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mIPV4\033[0m":
		wsIP4()
	case "2. \033[93mIPV6\033[0m":
		wsIP6()
	case "0. \033[94mBack to the main menu\033[0m":
	    clearScreen()
		mainMenu()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
func wsIP4() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Reverse \033[92mWS \033[96mTCP IPV4\033[0m ")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mIRAN\033[0m", "2. \033[93mKHAREJ\033[92m[1]\033[0m", "3. \033[93mKHAREJ\033[92m[2]\033[0m", "4. \033[93mKHAREJ\033[92m[3]\033[0m", "5. \033[93mKHAREJ\033[92m[4]\033[0m", "0. \033[94mBack to the previous menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mIRAN\033[0m":
		iranWst()
	case "2. \033[93mKHAREJ\033[92m[1]\033[0m":
		kharejWst()
	case "3. \033[93mKHAREJ\033[92m[2]\033[0m":
		kharejWst()
	case "4. \033[93mKHAREJ\033[92m[3]\033[0m":
		kharejWst()
	case "5. \033[93mKHAREJ\033[92m[4]\033[0m":
		kharejWst()
	case "0. \033[94mBack to the previous menu\033[0m":
	    clearScreen()
		wsMenu1()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
func wsIP6() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Reverse \033[92mWS \033[96mTCP IPV6\033[0m ")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mIRAN\033[0m", "2. \033[93mKHAREJ\033[92m[1]\033[0m", "3. \033[93mKHAREJ\033[92m[2]\033[0m", "4. \033[93mKHAREJ\033[92m[3]\033[0m", "5. \033[93mKHAREJ\033[92m[4]\033[0m", "0. \033[94mBack to the previous menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mIRAN\033[0m":
		iranWst()
	case "2. \033[93mKHAREJ\033[92m[1]\033[0m":
		kharejWst6()
	case "3. \033[93mKHAREJ\033[92m[2]\033[0m":
		kharejWst6()
    case "4. \033[93mKHAREJ\033[92m[3]\033[0m":
		kharejWst6()
	case "5. \033[93mKHAREJ\033[92m[4]\033[0m":
		kharejWst6()
	case "0. \033[94mBack to the previous menu\033[0m":
	    clearScreen()
		wsMenu1()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
const service1 = `[Unit]
Description=Rtun server

[Service]
Type=simple
ExecStart=/root/reverse-tunnel/./rtun-server -f /root/reverse-tunnel/rtun-server.yml
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
`
func iranWst() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m IRAN\033[92m TCP \033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	if _, err := os.Stat("/root/reverse-tunnel"); os.IsNotExist(err) {
		install1()
	} else {
		fmt.Println("\033[93mSkipping download..\033[0m")
	}

	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	displayNotification("Configuring Iran")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	var numServers int
	for {
		fmt.Print("\033[93mHow many \033[92mKharej Servers\033[93m do you have?\033[0m ")
		var numServersStr string
		fmt.Scanln(&numServersStr)
		numServersStr = strings.TrimSpace(numServersStr)
		var err error
		numServers, err = strconv.Atoi(numServersStr)
		if err == nil && numServers >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	var tunnelPort int
	for {
		fmt.Print("\033[93mEnter \033[92mTunnel port: \033[0m")
		var tunnelPortStr string
		fmt.Scanln(&tunnelPortStr)
		tunnelPortStr = strings.TrimSpace(tunnelPortStr)
		var err error
		tunnelPort, err = strconv.Atoi(tunnelPortStr)
		if err == nil && tunnelPort >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	agents := make([]string, numServers)
	for i := 0; i < numServers; i++ {
		authKey := genAuthKey()
		ports := genPorts(i + 1)

		serverInfo := fmt.Sprintf("- auth_key: %s\n  ports: %s", authKey, ports)
		agents[i] = serverInfo

		box := color.New(color.FgYellow)
		key := color.New(color.FgWhite, color.BgGreen)
		fmt.Println("╭──────────────────────────────────────────────────────────────╮")
		fmt.Printf("│                      %s                          │\n", box.Sprint("   Server Key ", i+1))
		fmt.Println("╰──────────────────────────────────────────────────────────────╯")
		fmt.Println("╭──────────────────────────────────────────────────────────────╮")
		fmt.Printf("%s\n", key.Sprint(authKey))
		fmt.Println("╰──────────────────────────────────────────────────────────────╯")
		fmt.Println()
	}

	configContent := fmt.Sprintf("control_address: 0.0.0.0:%d\n\nagents:\n%s", tunnelPort, strings.Join(agents, "\n"))
	err := ioutil.WriteFile("/root/reverse-tunnel/rtun-server.yml", []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("\033[91mfailed to write config:\033[0m %v\n", err)
		return
	}
	err = iranService()
	if err != nil {
		fmt.Printf("\033[91mCouldn't create iran service:\033[0m %v\n", err)
		return
	}
    resIran()
	fmt.Println("\033[92mIran Service Created Successfully.\033[0m")
}


func genPorts(serverNum int) string {
	fmt.Printf("\033[93mEnter \033[92mKharej \033[93m%d \033[96mPorts \033[93m(use comma between ports):\033[0m ", serverNum)
	var portsStr string
	fmt.Scanln(&portsStr)
	portsStr = strings.TrimSpace(portsStr)
	portsSlice := strings.Split(portsStr, ",")

	ports := make([]string, len(portsSlice))
	for i, port := range portsSlice {
		port = strings.TrimSpace(port)
		ports[i] = port + "/tcp"
	}

	return fmt.Sprintf("[%s]", strings.Join(ports, ", "))
}

func iranService() error {
	content := []byte(service1)

	err := ioutil.WriteFile("/etc/systemd/system/rtun-iran.service", content, 0644)
	if err != nil {
		return err
	}
	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		return fmt.Errorf("\033[91mFailed to reload:\033[0m %v", err)
	}
	err = exec.Command("systemctl", "enable", "rtun-iran").Run()
	if err != nil {
		return fmt.Errorf("\033[91mFailed to enable Service:\033[0m %v", err)
	}
	err = exec.Command("systemctl", "restart", "rtun-iran").Run()
	if err != nil {
		return fmt.Errorf("\033[91mFailed to restart Service!:\033[0m %v", err)
	}

	return nil
}
const service2 = `[Unit]
Description=Rtun client

[Service]
Type=simple
ExecStart=/root/reverse-tunnel/./rtun-client -f /root/reverse-tunnel/rtun-client.yml
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
`

func kharejWst() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Kharej\033[92m TCP IPV4 \033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	if _, err := os.Stat("/root/reverse-tunnel"); os.IsNotExist(err) {
		install2()
	} else {
		fmt.Println("\033[93mSkipping download..\033[0m")
	}

	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	displayNotification("Configuring Kharej")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	var iranIP string
	fmt.Print("\033[93mEnter \033[92mIran IPV4\033[93m address:\033[0m ")
	fmt.Scanln(&iranIP)
	iranIP = strings.TrimSpace(iranIP)

	var tunnelPort int
	for {
		fmt.Print("\033[93mEnter \033[92mTunnel port: \033[0m")
		var tunnelPortStr string
		fmt.Scanln(&tunnelPortStr)
		tunnelPortStr = strings.TrimSpace(tunnelPortStr)
		var err error
		tunnelPort, err = strconv.Atoi(tunnelPortStr)
		if err == nil && tunnelPort >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	var authKey string
	fmt.Print("\033[93mEnter \033[92mIRAN \033[93mkey:\033[0m ")
	fmt.Scanln(&authKey)
	authKey = strings.TrimSpace(authKey)

	var configPorts []string
	for {
		fmt.Print("\033[93mEnter \033[92mconfig ports \033[93m(use comma between ports):\033[0m ")
		var configPortsStr string
		fmt.Scanln(&configPortsStr)
		configPortsStr = strings.TrimSpace(configPortsStr)
		configPortsSlice := strings.Split(configPortsStr, ",")

		configPorts = make([]string, len(configPortsSlice))
		for i, port := range configPortsSlice {
			port = strings.TrimSpace(port)
			configPorts[i] = port + "/tcp"
		}

		if len(configPorts) >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	configContent := fmt.Sprintf("gateway_url: ws://%s:%d\nauth_key: %s\n\nforwards:\n", iranIP, tunnelPort, authKey)
	for _, port := range configPorts {
		configContent += fmt.Sprintf("  - port: %s\n    destination: 127.0.0.1:%s\n", port, strings.Split(port, "/")[0])
	}

	err := ioutil.WriteFile("/root/reverse-tunnel/rtun-client.yml", []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("\033[91mCouldn't write the config:\033[0m %v\n", err)
		return
	}

	err = KharejService()
	if err != nil {
		fmt.Printf("\033[91mCouldn't create kharej service:\033[0m %v\n", err)
		return
	}
    resKharej()
	fmt.Println("\033[92mKharej Service Created Successfully.\033[0m")
}

func KharejService() error {
	content := []byte(service2)

	err := ioutil.WriteFile("/etc/systemd/system/rtun-kharej.service", content, 0644)
	if err != nil {
		return err
	}
	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		return fmt.Errorf("\033[91mFailed to reload:\033[0m %v", err)
	}
	err = exec.Command("systemctl", "enable", "rtun-kharej").Run()
	if err != nil {
		return fmt.Errorf("\033[91mFailed to enable Service:\033[0m %v", err)
	}
	err = exec.Command("systemctl", "restart", "rtun-kharej").Run()
	if err != nil {
		return fmt.Errorf("\033[91mFailed to restart Service!:\033[0m %v", err)
	}

	return nil
}
func kharejWst6() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Kharej\033[92m TCP IPV6 \033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	if _, err := os.Stat("/root/reverse-tunnel"); os.IsNotExist(err) {
		install2()
	} else {
		fmt.Println("\033[93mSkipping download..\033[0m")
	}

	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	displayNotification("Configuring Kharej")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	var iranIP string
	fmt.Print("\033[93mEnter \033[92mIran IPV6\033[93m address:\033[0m ")
	fmt.Scanln(&iranIP)
	iranIP = strings.TrimSpace(iranIP)

	var tunnelPort int
	for {
		fmt.Print("\033[93mEnter \033[92mTunnel port: \033[0m")
		var tunnelPortStr string
		fmt.Scanln(&tunnelPortStr)
		tunnelPortStr = strings.TrimSpace(tunnelPortStr)
		var err error
		tunnelPort, err = strconv.Atoi(tunnelPortStr)
		if err == nil && tunnelPort >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	var authKey string
	fmt.Print("\033[93mEnter \033[92mIRAN \033[93mkey:\033[0m ")
	fmt.Scanln(&authKey)
	authKey = strings.TrimSpace(authKey)

	var configPorts []string
	for {
		fmt.Print("\033[93mEnter \033[92mconfig ports \033[93m(use comma between ports):\033[0m ")
		var configPortsStr string
		fmt.Scanln(&configPortsStr)
		configPortsStr = strings.TrimSpace(configPortsStr)
		configPortsSlice := strings.Split(configPortsStr, ",")

		configPorts = make([]string, len(configPortsSlice))
		for i, port := range configPortsSlice {
			port = strings.TrimSpace(port)
			configPorts[i] = port + "/tcp"
		}

		if len(configPorts) >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	configContent := fmt.Sprintf("gateway_url: ws://[%s]:%d\nauth_key: %s\n\nforwards:\n", iranIP, tunnelPort, authKey)
	for _, port := range configPorts {
		configContent += fmt.Sprintf("  - port: %s\n    destination: 127.0.0.1:%s\n", port, strings.Split(port, "/")[0])
	}

	err := ioutil.WriteFile("/root/reverse-tunnel/rtun-client.yml", []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("\033[91mCouldn't write the config:\033[0m %v\n", err)
		return
	}

	err = KharejService()
	if err != nil {
		fmt.Printf("\033[91mCouldn't create kharej service:\033[0m %v\n", err)
		return
	}
    resKharej()
	fmt.Println("\033[92mKharej Service Created Successfully.\033[0m")
}
func wsMenu2() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Reverse \033[92mWS \033[96mUDP\033[0m ")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mIPV4\033[0m", "2. \033[93mIPV6\033[0m", "0. \033[94mBack to the main menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mIPV4\033[0m":
		wsIpu4()
	case "2. \033[93mIPV6\033[0m":
		wsIpu6()
	case "0. \033[94mBack to the main menu\033[0m":
	    clearScreen()
		mainMenu()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
func wsIpu4() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Reverse \033[92mWS \033[96mUDP IPV4\033[0m ")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mIRAN\033[0m", "2. \033[93mKHAREJ\033[92m[1]\033[0m", "3. \033[93mKHAREJ\033[92m[2]\033[0m", "4. \033[93mKHAREJ\033[92m[3]\033[0m", "5. \033[93mKHAREJ\033[92m[4]\033[0m", "0. \033[94mBack to the previous menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mIRAN\033[0m":
		iranWsu()
	case "2. \033[93mKHAREJ\033[92m[1]\033[0m":
		kharejWsu()
	case "3. \033[93mKHAREJ\033[92m[2]\033[0m":
		kharejWsu()
	case "4. \033[93mKHAREJ\033[92m[3]\033[0m":
		kharejWsu()
	case "5. \033[93mKHAREJ\033[92m[4]\033[0m":
		kharejWsu()
	case "0. \033[94mBack to the previous menu\033[0m":
	    clearScreen()
		wsMenu2()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
func wsIpu6() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Reverse \033[92mWS \033[96mUDP IPV6\033[0m ")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mIRAN\033[0m", "2. \033[93mKHAREJ\033[92m[1]\033[0m", "3. \033[93mKHAREJ\033[92m[2]\033[0m", "4. \033[93mKHAREJ\033[92m[3]\033[0m", "5. \033[93mKHAREJ\033[92m[4]\033[0m", "0. \033[94mBack to the previous menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mIRAN\033[0m":
		iranWsu()
	case "2. \033[93mKHAREJ\033[92m[1]\033[0m":
		kharejWsu6()
	case "3. \033[93mKHAREJ\033[92m[2]\033[0m":
		kharejWsu6()
	case "4. \033[93mKHAREJ\033[92m[3]\033[0m":
		kharejWsu6()
	case "5. \033[93mKHAREJ\033[92m[4]\033[0m":
		kharejWsu6()
	case "0. \033[94mBack to the previous menu\033[0m":
	    clearScreen()
		wsMenu2()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
func iranWsu() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m IRAN\033[92m UDP \033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	if _, err := os.Stat("/root/reverse-tunnel"); os.IsNotExist(err) {
		install1()
	} else {
		fmt.Println("\033[93mSkipping download..\033[0m")
	}

	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	displayNotification("Configuring Iran")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	var numServers int
	for {
		fmt.Print("\033[93mHow many \033[92mKharej Servers\033[93m do you have?\033[0m ")
		var numServersStr string
		fmt.Scanln(&numServersStr)
		numServersStr = strings.TrimSpace(numServersStr)
		var err error
		numServers, err = strconv.Atoi(numServersStr)
		if err == nil && numServers >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	var tunnelPort int
	for {
		fmt.Print("\033[93mEnter \033[92mTunnel port: \033[0m")
		var tunnelPortStr string
		fmt.Scanln(&tunnelPortStr)
		tunnelPortStr = strings.TrimSpace(tunnelPortStr)
		var err error
		tunnelPort, err = strconv.Atoi(tunnelPortStr)
		if err == nil && tunnelPort >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	agents := make([]string, numServers)
	for i := 0; i < numServers; i++ {
		authKey := genAuthKey()
		ports := genPortu(i + 1)

		serverInfo := fmt.Sprintf("- auth_key: %s\n  ports: %s", authKey, ports)
		agents[i] = serverInfo

		box := color.New(color.FgYellow)
		key := color.New(color.FgWhite, color.BgGreen)
		fmt.Println("╭──────────────────────────────────────────────────────────────╮")
		fmt.Printf("│                      %s                         │\n", box.Sprint("   Server Key ", i+1))
		fmt.Println("╰──────────────────────────────────────────────────────────────╯")
		fmt.Println("╭──────────────────────────────────────────────────────────────╮")
		fmt.Printf("%s\n", key.Sprint(authKey))
		fmt.Println("╰──────────────────────────────────────────────────────────────╯")
		fmt.Println()
	}

	configContent := fmt.Sprintf("control_address: 0.0.0.0:%d\n\nagents:\n%s", tunnelPort, strings.Join(agents, "\n"))
	err := ioutil.WriteFile("/root/reverse-tunnel/rtun-server.yml", []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("\033[91mfailed to write config:\033[0m %v\n", err)
		return
	}
	err = iranService()
	if err != nil {
		fmt.Printf("\033[91mCouldn't create iran service:\033[0m %v\n", err)
		return
	}
    resIran()
	fmt.Println("\033[92mIran Service Created Successfully.\033[0m")
}
func genPortu(serverNum int) string {
	fmt.Printf("\033[93mEnter \033[92mKharej \033[93m%d \033[96mPorts \033[93m(use comma between ports):\033[0m ", serverNum)
	var portsStr string
	fmt.Scanln(&portsStr)
	portsStr = strings.TrimSpace(portsStr)
	portsSlice := strings.Split(portsStr, ",")

	ports := make([]string, len(portsSlice))
	for i, port := range portsSlice {
		port = strings.TrimSpace(port)
		ports[i] = port + "/udp"
	}

	return fmt.Sprintf("[%s]", strings.Join(ports, ", "))
}
func kharejWsu() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Kharej\033[92m UDP IPV4 \033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	if _, err := os.Stat("/root/reverse-tunnel"); os.IsNotExist(err) {
		install2()
	} else {
		fmt.Println("\033[93mSkipping download..\033[0m")
	}

	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	displayNotification("Configuring Kharej")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	var iranIP string
	fmt.Print("\033[93mEnter \033[92mIran IPV4\033[93m address:\033[0m ")
	fmt.Scanln(&iranIP)
	iranIP = strings.TrimSpace(iranIP)

	var tunnelPort int
	for {
		fmt.Print("\033[93mEnter \033[92mTunnel port: \033[0m")
		var tunnelPortStr string
		fmt.Scanln(&tunnelPortStr)
		tunnelPortStr = strings.TrimSpace(tunnelPortStr)
		var err error
		tunnelPort, err = strconv.Atoi(tunnelPortStr)
		if err == nil && tunnelPort >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	var authKey string
	fmt.Print("\033[93mEnter \033[92mIRAN \033[93mkey:\033[0m ")
	fmt.Scanln(&authKey)
	authKey = strings.TrimSpace(authKey)

	var configPorts []string
	for {
		fmt.Print("\033[93mEnter \033[92mconfig ports \033[93m(use comma between ports):\033[0m ")
		var configPortsStr string
		fmt.Scanln(&configPortsStr)
		configPortsStr = strings.TrimSpace(configPortsStr)
		configPortsSlice := strings.Split(configPortsStr, ",")

		configPorts = make([]string, len(configPortsSlice))
		for i, port := range configPortsSlice {
			port = strings.TrimSpace(port)
			configPorts[i] = port + "/udp"
		}

		if len(configPorts) >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	configContent := fmt.Sprintf("gateway_url: ws://%s:%d\nauth_key: %s\n\nforwards:\n", iranIP, tunnelPort, authKey)
	for _, port := range configPorts {
		configContent += fmt.Sprintf("  - port: %s\n    destination: 127.0.0.1:%s\n", port, strings.Split(port, "/")[0])
	}

	err := ioutil.WriteFile("/root/reverse-tunnel/rtun-client.yml", []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("\033[91mCouldn't write the config:\033[0m %v\n", err)
		return
	}

	err = KharejService()
	if err != nil {
		fmt.Printf("\033[91mCouldn't create kharej service:\033[0m %v\n", err)
		return
	}
    resKharej()
	fmt.Println("\033[92mKharej Service Created Successfully.\033[0m")
}
func kharejWsu6() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Kharej\033[92m UDP IPV6 \033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	if _, err := os.Stat("/root/reverse-tunnel"); os.IsNotExist(err) {
		install2()
	} else {
		fmt.Println("\033[93mSkipping download..\033[0m")
	}

	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	displayNotification("Configuring Kharej")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	var iranIP string
	fmt.Print("\033[93mEnter \033[92mIran IPV6\033[93m address:\033[0m ")
	fmt.Scanln(&iranIP)
	iranIP = strings.TrimSpace(iranIP)

	var tunnelPort int
	for {
		fmt.Print("\033[93mEnter \033[92mTunnel port: \033[0m")
		var tunnelPortStr string
		fmt.Scanln(&tunnelPortStr)
		tunnelPortStr = strings.TrimSpace(tunnelPortStr)
		var err error
		tunnelPort, err = strconv.Atoi(tunnelPortStr)
		if err == nil && tunnelPort >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	var authKey string
	fmt.Print("\033[93mEnter \033[92mIRAN \033[93mkey:\033[0m ")
	fmt.Scanln(&authKey)
	authKey = strings.TrimSpace(authKey)

	var configPorts []string
	for {
		fmt.Print("\033[93mEnter \033[92mconfig ports \033[93m(use comma between ports):\033[0m ")
		var configPortsStr string
		fmt.Scanln(&configPortsStr)
		configPortsStr = strings.TrimSpace(configPortsStr)
		configPortsSlice := strings.Split(configPortsStr, ",")

		configPorts = make([]string, len(configPortsSlice))
		for i, port := range configPortsSlice {
			port = strings.TrimSpace(port)
			configPorts[i] = port + "/udp"
		}

		if len(configPorts) >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	configContent := fmt.Sprintf("gateway_url: ws://[%s]:%d\nauth_key: %s\n\nforwards:\n", iranIP, tunnelPort, authKey)
	for _, port := range configPorts {
		configContent += fmt.Sprintf("  - port: %s\n    destination: 127.0.0.1:%s\n", port, strings.Split(port, "/")[0])
	}

	err := ioutil.WriteFile("/root/reverse-tunnel/rtun-client.yml", []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("\033[91mCouldn't write the config:\033[0m %v\n", err)
		return
	}

	err = KharejService()
	if err != nil {
		fmt.Printf("\033[91mCouldn't create kharej service:\033[0m %v\n", err)
		return
	}
    resKharej()
	fmt.Println("\033[92mKharej Service Created Successfully.\033[0m")
}
func cmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		errorMsg := fmt.Sprintf("\033[31mfailed to run '%s %s': %v\033[0m", name, strings.Join(args, " "), err)
		return fmt.Errorf(errorMsg)
	}

	return nil
}

func hidden(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Run()
	if err != nil {
		errorMsg := fmt.Sprintf("\033[31mfailed to run '%s %s': %v\033[0m", name, strings.Join(args, " "), err)
		return fmt.Errorf(errorMsg)
	}

	return nil
}
func input(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)

	return input, nil
}
func acme() {
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	displayNotification("Installing acme & Getting Cert")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	err := requirments()
	if err != nil {
		log.Println("\033[91mfailed to install requirements:\033[0m", err)
	}

	err = dnlAcme()
	if err != nil {
		log.Println("\033[91mfailed to download acme:\033[0m", err)
	}
    fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	email, err := input("\033[93mEnter your \033[92mEmail\033[93m address:\033[0m ")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	if err != nil {
		log.Println("\033[91mThere was sth wrong with your email address:\033[0m", err)
	}

	err = registerAcc(email)
	if err != nil {
		log.Println("\033[91mCouldn't register account:\033[0m", err)
	}
    fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	subdomain, err := input("\033[93mEnter your \033[92mSubdomain\033[93m:\033[0m ")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	if err != nil {
		log.Println("\033[91mThere was sth wrong with your subdomain:\033[0m", err)
	}
	
	err = cert(subdomain)
	if err != nil {
		log.Println("\033[91mfailed to issue certificate:\033[0m", err)
	}

	displayCheckmark("\033[92mGetting Cert was Successful!\033[0m")
}

func requirments() error {
	err := hidden("apt", "install", "curl", "-y")
	if err != nil {
		return err
	}
	err = hidden("apt", "install", "socat", "-y")
	if err != nil {
		return err
	}

	return nil
}

func dnlAcme() error {
	err := cmd("curl", "-L", "-o", "acme.sh.tar.gz", "https://github.com/acmesh-official/acme.sh/archive/master.tar.gz")
	if err != nil {
		return err
	}
	err = cmd("tar", "xzf", "acme.sh.tar.gz")
	if err != nil {
		return err
	}
	err = os.Chdir("acme.sh-master")
	if err != nil {
		return err
	}

	err = cmd("sh", "acme.sh", "--install")
	if err != nil {
		return err
	}

	return nil
}

func registerAcc(email string) error {
	err := os.Setenv("PATH", fmt.Sprintf("%s:%s", os.Getenv("PATH"), "/root/.acme.sh"))
	if err != nil {
		return err
	}
	err = cmd("acme.sh", "--register-account", "--accountemail", email)
	if err != nil {
		return err
	}

	return nil
}

func cert(subdomain string) error {
	err := os.Setenv("PATH", fmt.Sprintf("%s:%s", os.Getenv("PATH"), "/root/.acme.sh"))
	if err != nil {
		return err
	}

	err = cmd("acme.sh", "--issue", "--standalone", "-d", subdomain)
	if err != nil {
		return err
	}

	return nil
}
func wsMenu3() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Reverse \033[92mWSS \033[96mTCP\033[0m ")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mIRAN\033[0m", "2. \033[93mKHAREJ\033[92m[1]\033[0m", "3. \033[93mKHAREJ\033[92m[2]\033[0m", "4. \033[93mKHAREJ\033[92m[3]\033[0m", "5. \033[93mKHAREJ\033[92m[4]\033[0m", "0. \033[94mBack to the main menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mIRAN\033[0m":
		iranWsst()
	case "2. \033[93mKHAREJ\033[92m[1]\033[0m":
		kharejWsst()
	case "3. \033[93mKHAREJ\033[92m[2]\033[0m":
		kharejWsst()
	case "4. \033[93mKHAREJ\033[92m[3]\033[0m":
		kharejWsst()
	case "5. \033[93mKHAREJ\033[92m[4]\033[0m":
		kharejWsst()
	case "0. \033[94mBack to the main menu\033[0m":
	    clearScreen()
		mainMenu()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
func iranWsst() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m IRAN WSS\033[92m TCP \033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\033[93mHave you already obtained a cert for Iran? (\033[92myes/\033[91mno)\033[93m:\033[0m ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	answer = strings.ToLower(strings.TrimSpace(answer))

	if answer == "yes" || answer == "y" {
		fmt.Println("\033[91mlemme skip real quick..\033[0m")
		goto SkipCert
	}

	acme()
SkipCert:
	if _, err := os.Stat("/root/reverse-tunnel"); os.IsNotExist(err) {
		install1()
	} else {
		fmt.Println("\033[93m───────────────────────────────────────\033[0m")
		fmt.Println("\033[93mSkipping download..\033[0m")
	}


	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	displayNotification("Configuring IRAN")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	var numServers int
	for {
		fmt.Print("\033[93mEnter \033[92mthe number\033[93m of \033[96mkharej servers\033[93m:\033[0m ")
		numServersStr, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		numServersStr = strings.TrimSpace(numServersStr)
		numServers, err = strconv.Atoi(numServersStr)
		if err == nil {
			break 
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input!\033[0m")
	}
	fmt.Print("\033[93mEnter \033[92mTunnel port\033[93m:\033[0m ")
	tunnelPort, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	tunnelPort = strings.TrimSpace(tunnelPort)
	fmt.Print("\033[93mEnter your \033[92mSub domain\033[93m:\033[0m ")
	domain, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	domain = strings.TrimSpace(domain)
	config := fmt.Sprintf("control_address: :%s\n\nlets_encrypt:\n  domain: %s\n\nagents:\n", tunnelPort, domain)
	authKeys := genAuthK(numServers)
	for i := 0; i < numServers; i++ {
		ports := genPortzs(i + 1)
		serverInfo := fmt.Sprintf("  - auth_key: %s\n    ports: [%s]\n", authKeys[i], ports)
		config += serverInfo
		box := color.New(color.FgYellow)
		key := color.New(color.FgWhite, color.BgGreen)
		fmt.Println("╭──────────────────────────────────────────────────────────────╮")
		fmt.Printf("│                      %s                         │\n", box.Sprint("   Server Key ", i+1))
		fmt.Println("╰──────────────────────────────────────────────────────────────╯")
		fmt.Println("╭──────────────────────────────────────────────────────────────╮")
		fmt.Printf("%s\n", key.Sprint(authKeys[i]))
		fmt.Println("╰──────────────────────────────────────────────────────────────╯")
		fmt.Println()
	}
	err = os.WriteFile("/root/reverse-tunnel/rtun-server.yml", []byte(config), 0644)
	if err != nil {
		fmt.Printf("\033[91mFailed to write the config: %v\n\033[0m", err)
		return
	}
	err = iranService()
	if err != nil {
		fmt.Printf("\033[91mCouldn't create iran service:\033[0m %v\n", err)
		return
	}
	resIran()
	fmt.Println("\033[92mIran Service Created Successfully.\033[0m")
}

func genAuthK(numKeys int) []string {
	authKeys := make([]string, numKeys)
	for i := 0; i < numKeys; i++ {
		authKeys[i] = genAuthKey()
	}
	return authKeys
}

func genAuthKey() string {
	cmd := exec.Command("openssl", "rand", "-hex", "32")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("\033[91mCouldn't generate the key:\033[0m %v\n", err)
		os.Exit(1)
	}

	return strings.TrimSpace(string(output))
}

func genPortzs(serverNum int) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\033[93mEnter \033[92mKharej %d\033[96m Config ports \033[93m(use comma between ports):\033[0m ", serverNum)
	portsStr, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	portsStr = strings.TrimSpace(portsStr)
	ports := strings.Split(portsStr, ",")
	formatPorts := make([]string, len(ports))
	for i, port := range ports {
		formatPorts[i] = strings.TrimSpace(fmt.Sprintf("%s/tcp", port))
	}

	return strings.Join(formatPorts, ",")
}
func kharejWsst() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Kharej\033[92m TCP \033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	if _, err := os.Stat("/root/reverse-tunnel"); os.IsNotExist(err) {
		install2()
	} else {
		fmt.Println("\033[93mSkipping download..\033[0m")
	}
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	var iranDomain string
	fmt.Print("\033[93mEnter \033[92mIran subdomain\033[93m:\033[0m ")
	fmt.Scanln(&iranDomain)
	iranDomain = strings.TrimSpace(iranDomain)

	var authKey string
	fmt.Print("\033[93mEnter \033[92miran key\033[93m:\033[0m ")
	fmt.Scanln(&authKey)
	authKey = strings.TrimSpace(authKey)
	var numConfigs int
	for {
		fmt.Print("\033[93mEnter the \033[92mnumber\033[93m of\033[96m configs:\033[0m ")
		var numConfigsStr string
		fmt.Scanln(&numConfigsStr)
		numConfigsStr = strings.TrimSpace(numConfigsStr)
		var err error
		numConfigs, err = strconv.Atoi(numConfigsStr)
		if err == nil && numConfigs >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	configContent := fmt.Sprintf("gateway_url: wss://%s\nauth_key: %s\n\nforwards:\n", iranDomain, authKey)
	for i := 0; i < numConfigs; i++ {
		var configPort string
		fmt.Printf("\033[93mEnter \033[92mconfig port %d: \033[0m", i+1)
		fmt.Scanln(&configPort)
		configPort = strings.TrimSpace(configPort)

		configContent += fmt.Sprintf("  - port: %s/tcp\n    destination: 127.0.0.1:%s\n", configPort, configPort)
	}

	err := ioutil.WriteFile("/root/reverse-tunnel/rtun-client.yml", []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("\033[91mCouldn't write the config:\033[0m %v\n", err)
		return
	}

	err = KharejService()
	if err != nil {
		fmt.Printf("\033[91mCouldn't create kharej service:\033[0m %v\n", err)
		return
	}
    resKharej()
	fmt.Println("\033[92mKharej Service Created Successfully.\033[0m")
}


func wsMenu4() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Reverse \033[92mWSS \033[96mUDP\033[0m ")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")

	prompt := &survey.Select{
		Message: "Enter your choice Please:",
		Options: []string{"1. \033[92mIRAN\033[0m", "2. \033[93mKHAREJ\033[92m[1]\033[0m", "3. \033[93mKHAREJ\033[92m[2]\033[0m", "4. \033[93mKHAREJ\033[92m[3]\033[0m", "5. \033[93mKHAREJ\033[92m[4]\033[0m", "0. \033[94mBack to the main menu\033[0m"},
	}
    
	var choice string
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		log.Fatalf("\033[91mCan't read user input, sry!:\033[0m %v", err)
	}

	switch choice {
	case "1. \033[92mIRAN\033[0m":
		iranWssu()
	case "2. \033[93mKHAREJ\033[92m[1]\033[0m":
		kharejWssu()
	case "3. \033[93mKHAREJ\033[92m[2]\033[0m":
		kharejWssu()
	case "4. \033[93mKHAREJ\033[92m[3]\033[0m":
		kharejWssu()
	case "5. \033[93mKHAREJ\033[92m[4]\033[0m":
		kharejWssu()
	case "0. \033[94mBack to the main menu\033[0m":
	    clearScreen()
		mainMenu()
	default:
		fmt.Println("\033[91mInvalid choice\033[0m")
	}

	readInput()
}
func iranWssu() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m IRAN WSS\033[92m UDP \033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\033[93mHave you already obtained a cert for Iran? (\033[92myes/\033[91mno)\033[93m:\033[0m ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	answer = strings.ToLower(strings.TrimSpace(answer))

	if answer == "yes" || answer == "y" {
		fmt.Println("\033[91mlemme skip real quick..\033[0m")
		goto SkipCert
	}

	acme()
SkipCert:
	if _, err := os.Stat("/root/reverse-tunnel"); os.IsNotExist(err) {
		install1()
	} else {
		fmt.Println("\033[93m───────────────────────────────────────\033[0m")
		fmt.Println("\033[93mSkipping download..\033[0m")
	}


	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	displayNotification("Configuring IRAN")
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	var numServers int
	for {
		fmt.Print("\033[93mEnter \033[92mthe number\033[93m of \033[96mkharej servers\033[93m:\033[0m ")
		numServersStr, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		numServersStr = strings.TrimSpace(numServersStr)
		numServers, err = strconv.Atoi(numServersStr)
		if err == nil {
			break 
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input!\033[0m")
	}

	fmt.Print("\033[93mEnter \033[92mTunnel port\033[93m:\033[0m ")
	tunnelPort, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	tunnelPort = strings.TrimSpace(tunnelPort)

	fmt.Print("\033[93mEnter your \033[92mSub domain\033[93m:\033[0m ")
	domain, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	domain = strings.TrimSpace(domain)

	config := fmt.Sprintf("control_address: :%s\n\nlets_encrypt:\n  domain: %s\n\nagents:\n", tunnelPort, domain)
	authKeys := genAuthK(numServers)
	for i := 0; i < numServers; i++ {
		ports := genPortzu(i + 1)
		serverInfo := fmt.Sprintf("  - auth_key: %s\n    ports: [%s]\n", authKeys[i], ports)
		config += serverInfo

		box := color.New(color.FgYellow)
		key := color.New(color.FgWhite, color.BgGreen)
		fmt.Println("╭──────────────────────────────────────────────────────────────╮")
		fmt.Printf("│                      %s                         │\n", box.Sprint("   Server Key ", i+1))
		fmt.Println("╰──────────────────────────────────────────────────────────────╯")
		fmt.Println("╭──────────────────────────────────────────────────────────────╮")
		fmt.Printf("%s\n", key.Sprint(authKeys[i]))
		fmt.Println("╰──────────────────────────────────────────────────────────────╯")
		fmt.Println()
	}

	err = os.WriteFile("/root/reverse-tunnel/rtun-server.yml", []byte(config), 0644)
	if err != nil {
		fmt.Printf("\033[91mFailed to write configuration file: %v\n\033[0m", err)
		return
	}
	
	err = iranService()
	if err != nil {
		fmt.Printf("\033[91mCouldn't create iran service:\033[0m %v\n", err)
		return
	}
    resIran()
	fmt.Println("\033[92mIran Service Created Successfully.\033[0m")
}

func genPortzu(serverNum int) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\033[93mEnter \033[92mKharej %d\033[96m Config ports \033[93m(use comma between ports):\033[0m ", serverNum)
	portsStr, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	portsStr = strings.TrimSpace(portsStr)
	ports := strings.Split(portsStr, ",")
	formatPorts := make([]string, len(ports))
	for i, port := range ports {
		formatPorts[i] = strings.TrimSpace(fmt.Sprintf("%s/udp", port))
	}

	return strings.Join(formatPorts, ",")
}
func kharejWssu() {
	clearScreen()
	fmt.Println("\033[92m ^ ^\033[0m")
	fmt.Println("\033[92m(\033[91mO,O\033[92m)\033[0m")
	fmt.Println("\033[92m(   ) \033[93m Kharej\033[92m UDP \033[0m")
	fmt.Println("\033[92m \"-\" \033[93m════════════════════════════════════\033[0m")
	if _, err := os.Stat("/root/reverse-tunnel"); os.IsNotExist(err) {
		install2()
	} else {
		fmt.Println("\033[93mSkipping download..\033[0m")
	}
	fmt.Println("\033[93m───────────────────────────────────────\033[0m")
	var iranDomain string
	fmt.Print("\033[93mEnter \033[92mIran subdomain\033[93m:\033[0m ")
	fmt.Scanln(&iranDomain)
	iranDomain = strings.TrimSpace(iranDomain)

	var authKey string
	fmt.Print("\033[93mEnter \033[92miran key\033[93m:\033[0m ")
	fmt.Scanln(&authKey)
	authKey = strings.TrimSpace(authKey)
	var numConfigs int
	for {
		fmt.Print("\033[93mEnter the \033[92mnumber\033[93m of\033[96m configs:\033[0m ")
		var numConfigsStr string
		fmt.Scanln(&numConfigsStr)
		numConfigsStr = strings.TrimSpace(numConfigsStr)
		var err error
		numConfigs, err = strconv.Atoi(numConfigsStr)
		if err == nil && numConfigs >= 1 {
			break
		}
		fmt.Println("\033[91mInvalid input. Plz enter a valid input.\033[0m")
	}

	configContent := fmt.Sprintf("gateway_url: wss://%s\nauth_key: %s\n\nforwards:\n", iranDomain, authKey)
	for i := 0; i < numConfigs; i++ {
		var configPort string
		fmt.Printf("\033[93mEnter \033[92mconfig port %d: \033[0m", i+1)
		fmt.Scanln(&configPort)
		configPort = strings.TrimSpace(configPort)

		configContent += fmt.Sprintf("  - port: %s/udp\n    destination: 127.0.0.1:%s\n", configPort, configPort)
	}

	err := ioutil.WriteFile("/root/reverse-tunnel/rtun-client.yml", []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("\033[91mCouldn't write the config:\033[0m %v\n", err)
		return
	}

	err = KharejService()
	if err != nil {
		fmt.Printf("\033[91mCouldn't create kharej service:\033[0m %v\n", err)
		return
	}
    resKharej()
	fmt.Println("\033[92mKharej Service Created Successfully.\033[0m")
}
