#! /bin/bash
docker help &> /dev/null

if [[ "${?}" = 0 ]]; then
    docker build --no-cache -t globo/paredao:latest -f backend/Dockerfile backend/
    docker build --no-cache -t globo/nginx:1.25 -f infra/Dockerfile infra
else
    echo "É necessário ter o docker instalado"
fi

cd infra/
docker compose up