# 设置基础 URL
$base = "http://get.mereith.com/tagger/"

# 检查架构
$arch = if ([Environment]::Is64BitOperatingSystem) {
    if ([System.Runtime.InteropServices.RuntimeInformation]::ProcessArchitecture -eq [System.Runtime.InteropServices.Architecture]::Arm64) {
        "arm64"
    } else {
        "amd64"
    }
} else {
    Write-Host "不支持的架构"
    exit 1
}

$tagger = "tagger_windows_${arch}.exe"

Write-Host "系统架构为: $arch"
Write-Host "正在下载 $tagger ..."

# 创建临时目录
$tempPath = Join-Path $env:TEMP "tagger_temp"
New-Item -ItemType Directory -Force -Path $tempPath | Out-Null

# 下载文件
$downloadPath = Join-Path $tempPath "tagger.exe"
try {
    Invoke-WebRequest -Uri "$base$tagger" -OutFile $downloadPath
} catch {
    Write-Host "下载失败: $_"
    exit 1
}

# 创建安装目录
$installDir = Join-Path $env:LOCALAPPDATA "Tagger"
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

# 删除旧版本
if (Test-Path (Join-Path $installDir "tagger.exe")) {
    Remove-Item (Join-Path $installDir "tagger.exe") -Force
}

# 复制新版本
try {
    Copy-Item $downloadPath (Join-Path $installDir "tagger.exe") -Force
    
    # 添加到 PATH
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$installDir*") {
        [Environment]::SetEnvironmentVariable(
            "Path",
            "$userPath;$installDir",
            "User"
        )
        Write-Host "已添加到用户 PATH 环境变量"
    }
    
    Write-Host "安装完成！"
    Write-Host "安装位置: $installDir"
} catch {
    Write-Host "安装失败: $_"
    exit 1
} finally {
    # 清理临时文件
    Remove-Item -Path $tempPath -Recurse -Force
}
