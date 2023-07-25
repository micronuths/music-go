rm -rf "$(pwd)/service.file/music-web.service" /etc/systemd/system/music-web.service
mkdir -p "$(pwd)/service.file"
echo "\
[Unit]
Description=music-web
After=network.target
 
[Service]
Type=simple
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/music-web
Restart=on-failure
 
[Install]
WantedBy=multi-user.target\
" >> "$(pwd)/service.file/music-web.service"

cp $(pwd)/service.file/music-web.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable music-web.service