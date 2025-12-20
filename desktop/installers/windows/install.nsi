; DrCom 自动认证安装脚本
!include "MUI2.nsh"
!include "LogicLib.nsh"

; --- 基础配置 ---
!define APP_NAME "DrCom 自动认证"
!define EXE_NAME "drcom.exe"
!define TASK_NAME "DrCom_Auto_iHebut"
!define PUBLISHER "simplenty"
!define VERSION "1.0.0.0"

Name "${APP_NAME}"
OutFile "DrCom_Setup.exe"
InstallDir "$PROGRAMFILES\DrComClient"
InstallDirRegKey HKLM "Software\${APP_NAME}" ""
RequestExecutionLevel admin
Unicode true
ManifestDPIAware true

; --- 版本信息 ---
VIProductVersion "${VERSION}"
VIAddVersionKey "ProductName" "${APP_NAME}"
VIAddVersionKey "CompanyName" "${PUBLISHER}"
VIAddVersionKey "FileVersion" "${VERSION}"
VIAddVersionKey "FileDescription" "DrCom 自动认证安装程序"
VIAddVersionKey "LegalCopyright" "Copyright © 2025 ${PUBLISHER}"

; --- 界面流程 ---
!define MUI_ICON "..\..\assets\drcom.ico"
!define MUI_UNICON "..\..\assets\drcom.ico"
!define MUI_ABORTWARNING

!insertmacro MUI_PAGE_WELCOME
!define MUI_LICENSEPAGE_CHECKBOX
!insertmacro MUI_PAGE_LICENSE "..\..\..\LICENSE"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!define MUI_FINISHPAGE_RUN "$INSTDIR\${EXE_NAME}"
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_LANGUAGE "SimpChinese"

; --- 安装主程序 ---
Section "Install"
    SetOutPath "$INSTDIR"
    SetOverwrite on
    File "${EXE_NAME}"

    ; 修复目录权限 (防止配置文件无法写入)
    ExecWait 'icacls "$INSTDIR" /grant "Users":(OI)(CI)F /t /q'

    ; 写入卸载信息
    WriteUninstaller "$INSTDIR\Uninstall.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayName" "${APP_NAME}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "UninstallString" "$INSTDIR\Uninstall.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayIcon" "$INSTDIR\${EXE_NAME}"

    ; 创建开始菜单快捷方式 (所有用户)
    SetShellVarContext all
    CreateDirectory "$SMPROGRAMS\${APP_NAME}"
    CreateShortcut "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk" "$INSTDIR\${EXE_NAME}"
    CreateShortcut "$SMPROGRAMS\${APP_NAME}\卸载.lnk" "$INSTDIR\Uninstall.exe"

    ; 创建桌面快捷方式 (当前用户)
    SetShellVarContext current
    CreateShortcut "$DESKTOP\${APP_NAME}.lnk" "$INSTDIR\${EXE_NAME}"

    ; 注册计划任务
    Call CreateTaskXML
SectionEnd

; --- 卸载程序 ---
Section "Uninstall"
    ExecWait 'schtasks /delete /tn "${TASK_NAME}" /f'

    Delete "$INSTDIR\${EXE_NAME}"
    Delete "$INSTDIR\Uninstall.exe"
    Delete "$INSTDIR\config.json"
    RMDir "$INSTDIR"

    SetShellVarContext current
    Delete "$DESKTOP\${APP_NAME}.lnk"

    SetShellVarContext all
    RMDir /r "$SMPROGRAMS\${APP_NAME}"

    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}"
    DeleteRegKey HKLM "Software\${APP_NAME}"
SectionEnd

; --- 辅助函数：生成可读的 XML 配置 ---
Function CreateTaskXML
    StrCpy $R0 "$INSTDIR\task_temp.xml"
    FileOpen $0 "$R0" w

    ; 头部信息
    FileWrite $0 '<?xml version="1.0" encoding="UTF-16"?>$\r$\n'
    FileWrite $0 '<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">$\r$\n'
    FileWrite $0 '  <RegistrationInfo>$\r$\n'
    FileWrite $0 '    <Description>DrCom Auto Service</Description>$\r$\n'
    FileWrite $0 '  </RegistrationInfo>$\r$\n'

    ; 触发器: WLAN连接 iHebut 时触发
    FileWrite $0 '  <Triggers>$\r$\n'
    FileWrite $0 '    <EventTrigger>$\r$\n'
    FileWrite $0 '      <Enabled>true</Enabled>$\r$\n'
    FileWrite $0 '      <Subscription>&lt;QueryList&gt;&lt;Query Id="0" Path="Microsoft-Windows-WLAN-AutoConfig/Operational"&gt;&lt;Select Path="Microsoft-Windows-WLAN-AutoConfig/Operational"&gt;*[System[Provider[@Name=&apos;Microsoft-Windows-WLAN-AutoConfig&apos;] and (EventID=8001)]] and *[EventData[Data[@Name=&apos;SSID&apos;]=&apos;iHebut&apos;]]&lt;/Select&gt;&lt;/Query&gt;&lt;/QueryList&gt;</Subscription>$\r$\n'
    FileWrite $0 '    </EventTrigger>$\r$\n'
    FileWrite $0 '  </Triggers>$\r$\n'

    ; 运行权限
    FileWrite $0 '  <Principals>$\r$\n'
    FileWrite $0 '    <Principal id="Author">$\r$\n'
    FileWrite $0 '      <LogonType>InteractiveToken</LogonType>$\r$\n'
    FileWrite $0 '      <RunLevel>HighestAvailable</RunLevel>$\r$\n'
    FileWrite $0 '    </Principal>$\r$\n'
    FileWrite $0 '  </Principals>$\r$\n'

    ; 执行动作
    FileWrite $0 '  <Actions Context="Author">$\r$\n'
    FileWrite $0 '    <Exec>$\r$\n'
    FileWrite $0 '      <Command>"$INSTDIR\${EXE_NAME}"</Command>$\r$\n'
    FileWrite $0 '      <WorkingDirectory>$INSTDIR</WorkingDirectory>$\r$\n'
    FileWrite $0 '    </Exec>$\r$\n'
    FileWrite $0 '  </Actions>$\r$\n'

    ; 电源策略
    FileWrite $0 '  <Settings>$\r$\n'
    FileWrite $0 '    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>$\r$\n'
    FileWrite $0 '    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>$\r$\n'
    FileWrite $0 '    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>$\r$\n'
    FileWrite $0 '    <ExecutionTimeLimit>PT0S</ExecutionTimeLimit>$\r$\n'
    FileWrite $0 '  </Settings>$\r$\n'
    FileWrite $0 '</Task>$\r$\n'

    FileClose $0

    ; 导入并删除临时文件
    ExecWait 'schtasks /create /tn "${TASK_NAME}" /xml "$R0" /f'
    Delete "$R0"
FunctionEnd