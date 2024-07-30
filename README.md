# Votação BBB

Desenvolva um sistema de votação para o paredão do BBB, em versão web com HTML/CSS/Javascript e uma API REST como backend utilizando algumas das linguagem de programação abaixo
- Java
- Python
- Go
- Ruby

O paredão do BBB consiste em uma votação que confronta **dois integrantes** do programa BBB. A votação é apresentada em uma interface acessível pela web, na qual os usuários optam por votar em um dos integrantes apresentados. Uma vez realizado o voto, o usuário recebe uma tela com a confirmação do sucesso de seu voto e um panorama percentual dos votos por candidato até aquele momento. Além do frontend deve existir uma api em que de fato, será computado os votos.


Ideia inicial para o processamento em backend:

    Usar a linguagem Go para no backend, com uma fila para permitir o processamento assíncrono e um banco de dados (não?) relacional.

Requisitos:

- Testes unitários
- Testes de benchmarking em funções críticas