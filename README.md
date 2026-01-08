# Final Path - Estrutura do Projeto

## 📁 Estrutura de Arquivos

```
FinalPath/
├── main.go                      # Ponto de entrada do jogo
├── internal/
│   ├── config/
│   │   └── constants.go         # Constantes e configurações
│   ├── entity/
│   │   ├── enemy.go             # Lógica dos inimigos
│   │   ├── tower.go             # Lógica das torres
│   │   └── projectile.go        # Lógica dos projéteis
│   ├── game/
│   │   └── game.go              # Lógica principal do jogo
│   ├── gamemap/
│   │   └── map.go               # Lógica dos mapas e paths
│   └── utils/
│       └── utils.go             # Funções auxiliares
├── go.mod                       # Dependências do Go
└── README.md                    # Este arquivo
```

## 📄 Descrição dos Arquivos

### `main.go`
Ponto de entrada da aplicação. Contém apenas a função `main()` que:
- Cria uma nova instância do jogo
- Configura a janela
- Inicia o loop do jogo

### `internal/config/constants.go`
Define todas as constantes usadas no jogo:
- Tamanhos das entidades (Enemy, Tower, Projectile, Path)
- Configurações da janela (largura, altura, título)

### `internal/utils/utils.go`
Funções auxiliares reutilizáveis:
- `CenteredPosition`: Helper para trabalhar com posições centralizadas
- `CenterInPath()`: Calcula o centro de um path
- `GetCenterFromTopLeft()`: Converte coordenadas
- `Min()`, `Max()`: Funções matemáticas básicas

### `internal/entity/enemy.go`
Gerencia os inimigos do jogo:
- `Enemy`: Struct que representa um inimigo
- `NewEnemy()`: Cria um novo inimigo
- `IsAlive()`: Verifica se o inimigo está vivo
- `TakeDamage()`: Aplica dano ao inimigo
- `FollowPath()`: Faz o inimigo seguir o caminho

### `internal/entity/tower.go`
Gerencia as torres de defesa:
- `Tower`: Struct que representa uma torre
- `NewTower()`: Cria uma nova torre
- `IsEnemyInRange()`: Verifica se inimigo está no alcance
- `CanFire()`: Verifica se a torre pode atirar
- `Attack()`: Cria um projétil contra o inimigo
- `CanPlaceTower()`: Valida se pode colocar torre na posição

### `internal/entity/projectile.go`
Gerencia os projéteis:
- `Projectile`: Struct que representa um projétil
- `NewProjectile()`: Cria um novo projétil
- `Hit()`: Move o projétil e detecta colisão

### `internal/gamemap/map.go`
Gerencia os mapas e caminhos:
- `Path`: Struct que representa um segmento de caminho
- `Map`: Collection de paths
- `DefaultMap()`: Retorna o mapa padrão do jogo
- `Draw()`: Desenha o mapa na tela
- `IsPositionOnPath()`: Verifica se uma posição está no caminho

### `internal/game/game.go`
Gerencia o estado e lógica principal do jogo:
- `Game`: Struct principal que contém todo o estado
- `NewGame()`: Cria uma nova instância do jogo
- `Update()`: Atualiza o estado do jogo a cada frame
- `Draw()`: Renderiza tudo na tela
- `Layout()`: Define o tamanho da tela
- `handleTowerPlacement()`: Gerencia colocação de torres
- `drawBuildableAreas()`: Desenha áreas onde torres podem ser colocadas

## 🔧 Como Executar

```bash
# Compilar
go build

# Executar
./FinalPath

# Ou executar diretamente
go run .
```

## 📝 Convenções de Código

- Funções e tipos exportados começam com letra maiúscula (ex: `NewEnemy`)
- Funções e campos privados começam com letra minúscula (usados internamente)
- Pacotes no diretório `internal/` não podem ser importados por código externo
- Métodos de struct seguem o padrão Go (ex: `enemy.IsAlive()`)
- Todas as posições X/Y em entidades representam o **centro** da entidade
- Constantes são definidas em `internal/config/constants.go`
