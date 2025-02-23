#!/bin/bash
version="$1"
./ioncube_encoder.sh -C -x86-64 -$version filemahoa -o mahoadone
cd mahoadone
ZIP_NAME="mahoa_$(date +%Y%m%d%H%M%S).zip"
zip -r "$ZIP_NAME"
getServer() {
    response=$(curl -s "https://api.gofile.io/servers")
    echo $(jq -r '.data.servers[0].name' <<< "$response")
}

ZIP_FILE_PATH="${ZIP_NAME}"


SERVER=$(getServer)

if [ -z "$SERVER" ]; then
    echo "Không Thể Lấy Server Từ Gofile!"
    exit 1
fi

RESPONSE=$(curl -s -X POST "https://${SERVER}.gofile.io/contents/uploadfile" \
    -F "file=@${ZIP_FILE_PATH}")

DOWNLOAD_LINK=$(jq -r '.data.downloadPage' <<< "$RESPONSE")

if [ "$DOWNLOAD_LINK" != "null" ]; then
    cd ..
    rm -rf mahoadone
    echo "File Đã Được Mã Hóa Thành Công!"
    echo "Link Tải: $DOWNLOAD_LINK"
else
    echo "Upload Thất Bại!"
fi
