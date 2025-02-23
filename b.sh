#!/bin/bash

./ioncube_encoder.sh -C -x86-64 -81 filemahoa -o mahoadone
getServer() {
    response=$(curl -s "https://api.gofile.io/servers")
    echo $(jq -r '.data.servers[0].name' <<< "$response")
}

# Định nghĩa biến
TOKEN="<ĐƯỜNG_DẪN_TOKEN_CỦA_BẠN>"  # Nếu cần token
FILE_NAME="test.zip"  # Đổi thành tên file cần upload
ZIP_FILE_PATH="${TOKEN}/${FILE_NAME}"

# Lấy server khả dụng từ Gofile
SERVER=$(getServer)

# Kiểm tra nếu server không lấy được
if [ -z "$SERVER" ]; then
    echo "Không thể lấy server từ Gofile!"
    exit 1
fi

# Upload file lên Gofile
RESPONSE=$(curl -s -X POST "https://${SERVER}.gofile.io/contents/uploadfile" \
    -F "file=@${ZIP_FILE_PATH}")

# In phản hồi từ API
echo "Phản hồi từ Gofile: $RESPONSE"

# Lấy link tải file từ phản hồi API
DOWNLOAD_LINK=$(jq -r '.data.downloadPage' <<< "$RESPONSE")

if [ "$DOWNLOAD_LINK" != "null" ]; then
    echo "File đã được upload thành công!"
    echo "Link tải: $DOWNLOAD_LINK"
else
    echo "Upload thất bại!"
fi
