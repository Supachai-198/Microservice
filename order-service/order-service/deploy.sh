sshpass -p '$Admin_1jj395qu' ssh -o StrictHostKeyChecking=no root@159.223.86.11 << ENDSSH
    docker login -u $1 -p $2 $3
    docker stop $4 || true && docker rm $4 || true
    docker pull $5
    docker run --restart=always -d --name $4 -p 5020:6000 $5
ENDSSH