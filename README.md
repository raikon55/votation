# Votação BBB

Desenvolva um sistema de votação para o paredão do BBB, em versão web com HTML/CSS/Javascript e uma API REST como backend utilizando algumas das linguagem de programação abaixo
- Java
- Python
- Go
- Ruby

O paredão do BBB consiste em uma votação que confronta **dois integrantes** do programa BBB. A votação é apresentada em uma interface acessível pela web, na qual os usuários optam por votar em um dos integrantes apresentados. Uma vez realizado o voto, o usuário recebe uma tela com a confirmação do sucesso de seu voto e um panorama percentual dos votos por candidato até aquele momento. Além do frontend deve existir uma api em que de fato, será computado os votos.


Ideia inicial para o processamento em backend:

    Usar a linguagem Go para no backend, com uma fila para permitir o processamento assíncrono e um banco de dados (não?) relacional.

Ideia inicial para a visualização do frontend:

    Uma tela simples onde será exibido a foto dos dois atuais candidatos para o paredão, e com a possibilidade de votar em qualquer um deles ao clicar na foto.
    Retornar mensagem de sucesso. Exibir um percentual de votos no rodapé da página.

Ideia inicial para o infra:

    Um NGINX como proxy reverso e loadbalancer para permitir o alto fluxo de dados. Além disso o NGINX será responsável por criar uma primeira barreira para possíveis bot e DDoS.


---

Para um build do projeto, basta executar o script [build.sh](./build.sh). É necessário ter o docker instalado e ter acesso a um emulador de terminal Bash.

# Considerações sobre o projeto

Atualmente o projeto tem o backend escrito em Go, que é responsável por enviar cada voto para um fila de um serviço de mensagem, que nesse caso foi escolhido o RabbitMQ por simplicidade.

> O envio para fila garante permite um maior throughput da aplicação, pois tira o peso da abrir uma conexão com o banco, salvar o dado, aguarda uma confirmação do mesmo. Além disso, a própria fila serve de fallback para possíveis imprevistos q podem acarretar na perda de dados.

Após o voto ser registrado na fila, uma goroutine é encarregada de salvar no banco de forma apartada do serviço da API, além desse consumer pode ser desativado por variável de ambiente (ENABLE_CONSUMER), podendo escalar a API sem alocar novos consumidores para a fila. 

Por fim, há um NGINX como proxy reverso para expor apenas o backend e frontend, sem externalizar a fila e o banco de dados, e estes serão utilizados pela rede interna da stack. Isso facilita a migração entre contexto, como de docker compose para kubernetes, ou alguma arquitetura mais simples. O NGINX também é responsável por limitar as requisições.


# Trabalhos futuros

O próximos passos seria subir a aplicação em cluster K8S, e adicionar um kube-prometheus ao cluster, assim será possível criar visões e alertas baseados no Prometheus e no Grafana. Além disso, será implementado o SDK do Prometheus na própria aplicação para contabilizar as requisições, chamadas à fila e ao banco de dados, e possíveis erros.

Por fim, um frontend deve ser criado apenas com HTML, CSS e JS, para ser simples e pq ninguém merece frontend, com todo respeito.