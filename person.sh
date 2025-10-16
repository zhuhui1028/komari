env_check() {
    mach=$(uname -m)
    case "$mach" in
        amd64|x86_64)
            os_arch="amd64"
            ;;
        i386|i686)
            os_arch="386"
            ;;
        aarch64|arm64)
            os_arch="arm64"
            ;;
        *arm*)
            os_arch="arm"
            ;;
        s390x)
            os_arch="s390x"
            ;;
        riscv64)
            os_arch="riscv64"
            ;;
        mips)
            os_arch="mips"
            ;;
        mipsel|mipsle)
            os_arch="mipsle"
            ;;
        *)
            err "Unknown architecture: $uname"
            exit 1
            ;;
    esac

    system=$(uname)
    case "$system" in
        *Linux*)
            os="linux"
            ;;
        *Darwin*)
            os="darwin"
            ;;
        *FreeBSD*)
            os="freebsd"
            ;;
        *)
            err "Unknown architecture: $system"
            exit 1
            ;;
    esac
}
if [ $# -lt 1 ] ; then
echo "USAGE: $0 密钥"
exit 1;
fi

env_check
mkdir -p /opt/my-komari-agent
useradd --system --no-create-home --shell /bin/bash komari

cd /opt/my-komari-agent
version=1.1.12
komari_url="https://github.com/komari-monitor/komari-agent/releases/download/${version}/komari-agent-linux-${os_arch}"
wget $komari_url
mv komari-agent-linux-${os_arch} komari-agent

echo "ARGS=$1" > /opt/my-komari-agent/config.env
cat << 'EOF' > /etc/systemd/system/my-komari-agent.service
[Unit]
Description=Komari Agent Service
Documentation=other user
After=network.target

[Service]
# 指定运行服务的用户和用户组
User=komari
Group=komari

# 设置工作目录
WorkingDirectory=/opt/my-komari-agent

# 从环境文件中加载密钥
EnvironmentFile=/opt/my-komari-agent/config.env

# 启动服务的命令 (使用 ${KOMARI_TOKEN} 变量)
ExecStart=/opt/my-komari-agent/komari-agent -e https://komari.zhuhui.de -t ${ARGS}

# 设置失败时自动重启
Restart=always
RestartSec=10s

# 为日志指定一个标识符
SyslogIdentifier=my-komari-agent

[Install]
WantedBy=multi-user.target
EOF
chown -R komari:komari /opt/my-komari-agent
chmod 755 /opt/my-komari-agent/komari-agent

systemctl daemon-reload
systemctl start my-komari-agent.service
systemctl enable my-komari-agent.service
