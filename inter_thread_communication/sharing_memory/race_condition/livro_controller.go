package main

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func CadastraLivro(c *gin.Context) {
	var livro Livro
	c.ShouldBindJSON(&livro)

	var repository LivroRepository

	if repository.ExistsByIsbn(livro.Isbn) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"Erro": "Livro já existe"})
		return
	}

	repository.Save(livro)

	c.JSON(http.StatusCreated, livro)
}

/* Será que esse if garante que não vai ter livros duplicados na tabela?

Porque se um usuário cadastra um livro, e outro usuário cadastra o mesmo livro
ao mesmo tempo, antes do primeiro usuário receber a resposta, o segundo usuário
também vai cadastrar o livro, se os dois passarem pelo if da linha
15 ao mesmo tempo, o método vai retornar false para os dois

se rodarmos um 'select l.* from livro l', vai retornar os 2 registros
ou seja, o if não resolveu o problema, ele não impediu que o usuário
cadastrasse 2 livros iguais

esse é um bug de race condition, existe uma janela entre ler e inserir o dado
e nessa pequena janela de meio milisegundo, pode haver esse problema de concorrência

e se sua aplicação é web, ela naturalmente é concorrente.

para corrigir esse problema podemos configurar a coluna Isbn da tabela Livro
como "unique", porque assim o banco de dados garante que só será cadastrado
livros com ISBN únicos, perceba que os usuários ainda não vão entrar no if
mas só o primeiro vai conseguir cadastrar, o segundo vai dar erro ao tentar
cadastrar, ambos irão chamar o Save do repository.

precisamos nos preocupar com race condition em casos de transações financeiras e
atualizações de saldo, por exemplo, o que acontece se 1 cliente transferir R$100
e outro retirar R$50, ambos ao mesmo tempo?

podemos usar uma estrategia chamada mutex locking, que garante que apenas uma
thread por vez possa acessar o recurso protegido pelo mutex.
Se uma thread tentar obter o lock enquanto outro thread já o possui, a thread
que tenta obter o lock será colocada em espera (bloqueada) até que o lock seja
liberado pela thread que o possui. Isso efetivamente cria uma fila de threads,
onde cada thread aguarda sua vez para obter o lock e acessar o recurso compartilhado
*/

type Livro struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;"`
	Isbn   string
	Titulo string
	Autor  string
}

type LivroRepository struct{}

func (r *LivroRepository) ExistsByIsbn(isbn string) bool {
	//SELECT COUNT(*) AS total
	//FROM livros
	//WHERE isbn = isbn;
	return false
}

func (r *LivroRepository) Save(livro Livro) {
	// INSERT INTO livros (id, isbn, titulo, autor) VALUES.......
}
