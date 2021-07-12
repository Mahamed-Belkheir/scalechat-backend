docker build -f docker/Dockerfile.base -t scalechat_base:0.1 .;
docker build -f docker/Dockerfile.chat -t scalechat_chat:0.1 .;
docker build -f docker/Dockerfile.socket -t scalechat_socket:0.1 .;
docker build -f docker/Dockerfile.user -t scalechat_user:0.1 .;