# Rinha de Backend 2024 - Q1

Esta versão foi escrita em Golang usando Mux e Gorm. Sendo Mux um substituto popular ao clássico net/http para roteamento. E Gorm sendo uma biblioteca ORM viável para projetos em Golang. 

## Ferramentas:

- Golang: É a nossa linguagem queridinha, poderia ter escolhido outra opção mais preferi esperimentar na minha linguagem preferida.

- Gorm: Biblioteca de ORM capacitada a nós ajudar neste projeto.

- Mux: Biblioteca de roteamento  muito interessante seu uso. 

### Como gerar uma imagem local

Caso você deseje apenas esperimentar localmente a imagem e não seu correspondente docker.io você pode executar o seguinte commando 

>[!WARNING]
> Verifique se há espaço disponível em máquina


```bash
docker build -t rinha-backend/latest .
```

E o docker vai buildar sua imagem para testar o projeto localmente