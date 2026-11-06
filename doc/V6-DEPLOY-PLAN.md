# Deploy na v6 — plano de implementação

**Status:** implementado em 04/11/2026, falta validar contra a API real (seção 14) · **Data:** 03/11/2026 · **Autoria:** Patrick Selau Menoti

**Base:** guia do Console "Azion v3 · v4 · v6" (22/10/2026), SDK `github.com/aziontech/azionapi-v6-go-sdk` v0.2.0 e o código atual do deploy da v4 (`pkg/pipeline`, `pkg/cmd/deploy_remote`, `pkg/manifest`). Contexto e perguntas ao time de API: `doc/plan.md`, seções 13.1 e 13.2.

O deploy da v6 roda no mesmo `pipeline.Run` da v4, com uma lista de steps própria e um estado próprio. A diferença está dentro dos steps. Primeiro o CLI escreve a configuração em Drafts. Depois chama o build de cada versão e espera o estado final por polling. Por último monta uma release, faz o build dela, ativa e confere o resultado. Nada muda no tráfego até a ativação, então qualquer falha antes dela deixa o site como estava.

---

## 1. Escopo e premissas

| # | Premissa | Se estiver errada |
|---|---|---|
| P1 | O bundler gera um manifest da v6 que o CLI consegue ler (seção 8). O time do bundler ainda está trabalhando nisso. | O formato do manifest muda; os steps não. O contrato `ManifestV6` é o único ponto de ajuste. |
| P2 | Para contas v6, o CLI chama tudo em `/v6`, como o SDK define. | Se parte for `/v4`, o ajuste fica nos clientes de `pkg/api/v6` (URL por cliente no registry). Pergunta aberta em 13.2. |
| P3 | O CLI escreve sempre pelos paths com versão (`.../versions/{version_id}/...`), nunca pelos paths sem versão. | Nada quebra: os paths com versão são explícitos. Os sem versão têm semântica não documentada (13.2). |
| P4 | Uma conta resolve para uma geração por invocação (AS-2). A detecção por `use_v6_configurations` já existe. | Ver AS-2 em `doc/plan.md`. |

Fora deste plano: `init` e `link` para v6, gradual rollout (release `CANDIDATE`), promote entre Environments, rollback pelo CLI e o script do deploy remoto da v6 (D1, seção 12). Este plano cobre o pipeline local; o remoto terá um script próprio, novo ou adaptado do atual.

## 2. O que muda em relação ao deploy da v4

| Tema | v4 hoje | v6 |
|---|---|---|
| Efeito de salvar | Vai ao ar na hora. | Grava um Draft. Nada vai ao ar. |
| Como um step termina | Quando a chamada HTTP retorna. | Build e release disparam trabalho assíncrono; o step termina quando o estado chega a `ready`, `error` ou `canceled`. |
| Unidade que o CLI acompanha | id do recurso. | id do recurso **e** `version_id`. |
| Publicação | Implícita. | Release em um Deployment de um Environment: criar, fazer o build, ativar. |
| Falha no meio | Parte das mudanças já está no ar. | Nada muda no tráfego antes da ativação. |
| Recursos órfãos | Apagados no fim (`ManifestDeleteOrphanedResources`). | Dentro do Draft, sub-recursos órfãos são removidos, o que é seguro porque o Draft não está no ar. Recursos de primeiro nível não são apagados (seção 7). |
| Workload | Aponta direto para Application e Firewall. | Recebe um binding por Environment (Environment + Deployment + domínios). |

## 3. Fluxo do deploy

Os steps são uma `pipeline.Step[DeployStateV6]` cada, na ordem abaixo. Os nomes seguem a convenção do prefixo `Manifest` para o que aplica o manifest, e são contrato com o resumo de timing da v6 (seção 10).

| # | Step | Aplica quando | O que faz | Termina quando | Grava no azion.json |
|---|---|---|---|---|---|
| 1 | `GenerateManifest` | primeiro run sem manifest | Igual à v4. | — | — |
| 2 | `Build` | não `--skip-build` | Build do framework pelo bundler (P1). | — | prefix |
| 3 | `ReadManifest` | sempre | Lê e valida o `ManifestV6`. | Manifest válido. | — |
| 4 | `Bucket`, `Credentials`, `UploadStaticFiles` | manifest declara storage | Igual à v4. Storage não é versionado na v6. | Upload concluído. | bucket |
| 5 | `ManifestFunctions`, `ManifestNetworkLists`, `ManifestWAFs`, `ManifestConnectors` | manifest declara o recurso | Obtém um Draft (seção 5) e escreve a configuração nele. | Draft conferido (seção 6). | id, `version_id` |
| 6 | `VersionBuildDependencies` | algum Draft no step 5 | Chama o build de cada Draft do step 5 e espera todos (seção 4). | Todos `ready`, ou erro com a lista dos que falharam. | `version_state` |
| 7 | `ManifestApplication`, `ManifestFirewalls`, `ManifestCustomPages` | manifest declara o recurso | Obtém o Draft e escreve. Rules Engine, Cache Settings, Device Groups e Function Instances vão no Draft da Application. | Draft conferido. | id, `version_id` |
| 8 | `VersionBuildPublishable` | algum Draft no step 7 | Build e espera, como no step 6. | Todos `ready`. | `version_state` |
| 9 | `ResolveEnvironment` | sempre | Localiza o Environment pelo nome. O padrão é Production (D3); `--environment`, o azion.json ou o manifest podem apontar outro. Não cria Production, que já existe na conta. | Environment resolvido. | environment id |
| 10 | `ResolveDeployment` | sempre | Usa o Deployment gravado no azion.json ou cria um, com a política do Environment e `binding_policy` `FLEXIBLE` (como o bundler escreve). Num Deployment `STRICT`, a API recusa o build de toda release depois da primeira: 409 `43001` "The field 'resources' is immutable after creation", com a referência da Application lida como vazia (`got: ""`). | Deployment resolvido. | deployment id |
| 11 | `ReleaseCreate` | sempre | Cria a release em Draft apontando cada Application, Firewall e Custom Page, mais todas as dependências do manifest (Connectors, Functions, WAFs, Network Lists), para a versão que este run construiu (o `version_id` devolvido pela API, como `A32XCE2E`; D7). `origin.type` vai como `cli`. | Release em `draft`. | release id |
| 12 | `ReleaseBuild` | sempre | `POST .../build` e espera. | `ready`. Em `error`, a release está perdida e o próximo run cria outra. | release state |
| 13 | `ReleaseActivate` | não `--no-activate` | `POST .../activate` e espera o papel de tráfego. | `traffic_role` `ACTIVE` (single_version) ou `VALID_URL` (versioned_urls). | — |
| 14 | `ManifestWorkload` | manifest declara workload | Cria ou atualiza o Workload. Só mexe no binding deste Environment e preserva os outros. | Workload salvo. | workload id, url |
| 15 | `VerifyDeploy` | step 13 rodou | Confere o que está no ar (seção 6). | Conferência ok. | — |
| 16 | `ManifestPurge` | manifest declara purge | Igual à v4, depois da ativação. | — | — |

Por que os builds rodam em dois níveis (steps 6 e 8), e não todos de uma vez: o build de uma Application só termina quando os Connectors e Functions que ela usa estão `ready`. Se o build de um Connector falhar com todos disparados juntos, a Application fica presa ou falha com um erro que aponta para o lugar errado. Com dois níveis, o erro aparece no recurso que falhou.

A release aponta cada recurso para o `version_id` que este run construiu (D7), um short_id como `A32XCE2E`. A API não aceita `"latest"`: o campo exige um short_id de pelo menos 8 caracteres. Como os steps 6 e 8 só terminam com essas versões em `ready`, a release sempre nomeia versões prontas. A API não completa dependências: uma release sem um Connector que a Application usa falha no build com 400 `Dependency Mismatch` ("missing required subresource dependency of type: connector"). Como a versão também não lista do que depende, a release nomeia todas as dependências do manifest. Se o build da release é recusado, a release é apagada, como numa falha de build.

O risco é uma corrida: se o Console ou outro CLI terminar um build do mesmo recurso entre o step 8 e o build da release, é essa versão que vai ao ar. O CLI não impede isso, mas detecta: o `VerifyDeploy` compara a versão que entrou na release com a que este run construiu e avisa quando forem diferentes (seção 6).

## 4. Espera por estado (polling)

O SDK não tem webhook. Esperar é consultar o recurso até um estado final.

### 4.1 Estados

Um tipo só para recursos, Deployment, Environment e release, em `pkg/api/v6/version`:

```go
// State is the build state the v6 API reports as version_state (resources)
// or state (releases).
type State string

const (
	Draft    State = "draft"
	Queued   State = "queued"
	Building State = "building"
	Ready    State = "ready"
	Error    State = "error"
	Canceled State = "canceled"
	// archiving, archived, deleting, deleted
)

// Settled reports whether the state will not change without a new request.
func (s State) Settled() bool

// Failed reports whether a settled state means the build did not produce
// a usable version.
func (s State) Failed() bool
```

Um estado desconhecido não é tratado como final. O CLI continua esperando até o timeout e registra o valor no log de debug. Assim, um estado novo na API não vira sucesso nem falha por engano.

### 4.2 Helper de espera

Fica em um package pequeno e sem dependência da v6, `pkg/poll`, porque não tem nada de específico da v6:

```go
type Options struct {
	Interval time.Duration // intervalo entre consultas
	Timeout  time.Duration // limite total da espera
	// After substitui time.After nos testes, que rodam sem dormir.
	After func(time.Duration) <-chan time.Time
}

// Until calls fetch every Interval until done reports true, ctx ends or the
// timeout passes. A fetch error that Retryable accepts is retried on the next
// tick (or after Retry-After, when the API sends one); any other ends the wait.
func Until[T any](ctx context.Context, opt Options,
	fetch func(context.Context) (T, error),
	done func(T) (bool, error)) (T, error)
```

- Uma consulta a cada 5 segundos por recurso esperado (`Interval` = 5s, constante nomeada em `pkg/poll`), sem backoff (D5). O limite de taxa informado é de 200 requisições por segundo, a confirmar com o time de API. O pior caso do CLI fica muito abaixo disso: são sete tipos de recurso versionado, um nível de build por vez.
- Um `429` espera o `Retry-After` quando o header vier, e senão o próximo tick. Rede e `5xx` também são retentados no próximo tick. `4xx` encerra a espera com o erro da API.
- O timeout padrão é de 3 minutos por espera (`DefaultTimeout`, constante nomeada em `pkg/poll`), ou seja, por chamada de `Until`: cada build de recurso e o build da release têm o seu. `--timeout` substitui o padrão.
- O timeout devolve `poll.ErrTimeout` junto com o último estado visto. A mensagem diz qual recurso, qual versão e em que estado ficou.
- `ctx` cancelado (Ctrl-C) devolve `ctx.Err()` na hora (seção 7).

### 4.3 Builds de um nível

Os builds de um nível são disparados e esperados em paralelo, com `sync.WaitGroup` e um limite de concorrência. Uma falha **não** cancela os outros builds do mesmo nível: o step espera todos e devolve `errors.Join` com cada falha. O usuário vê de uma vez tudo o que deu errado. Os builds que deram certo viram versões Ready que ninguém referencia, o que é inofensivo.

Cada falha é um erro tipado:

```go
type BuildError struct {
	Kind      string // "application", "connector", ...
	Name      string
	ID        int64
	VersionID string
	State     version.State
	TraceID   string // quando a API devolver
}
```

### 4.4 Uma interface fina para recursos versionados

Os sete recursos com build (applications, connectors, custom_pages, firewalls, functions, network_lists, wafs) expõem as mesmas quatro operações. Uma interface documenta esse contrato, e os steps de build não repetem o mesmo código para cada recurso:

```go
type Versioned interface {
	CreateVersion(ctx context.Context, id int64, from, comment string) (string, error)
	Build(ctx context.Context, id int64, versionID, comment string) error
	CancelBuild(ctx context.Context, id int64, versionID string) error
	VersionState(ctx context.Context, id int64, versionID string) (version.State, error)
}
```

Cada cliente em `pkg/api/v6/<recurso>` implementa a interface. É o mesmo princípio do Tier 1: a interface só descreve o que os clientes já fazem.

## 5. Drafts: como o CLI obtém um Draft para escrever

| Situação no azion.json | Estado remoto | O que o CLI faz |
|---|---|---|
| Sem id | — | Cria o recurso. A resposta traz `version_id` em `draft`. |
| Com id e um `version_id` gravado | `draft`, `canceled` ou `error` | Reaproveita esse Draft: são estados editáveis. Evita acumular Drafts entre runs que falharam (limite padrão de 20 por recurso). |
| Com id, versão gravada `ready` ou ausente | — | `CreateVersion` sem `source_version`, que clona a última Ready, com comment `azion-cli deploy`. |
| Com id que não existe mais (404) | — | Cria de novo e avisa. Na v4 isso só acontece com cache settings e rules (`resource_applier.go`); na v6 vale para todo recurso versionado. |

A escrita é sempre completa: `PUT` na versão com a configuração do manifest, que é a fonte da verdade. Os sub-recursos do Draft (rules, cache settings, device groups, function instances, exceptions do WAF) são reconciliados por nome: atualiza quem existe, cria quem falta e remove do Draft quem não está mais no manifest. Remover dentro do Draft é seguro porque o Draft não está no ar.

## 6. Verificação

Três conferências, cada uma no ponto em que um erro ainda é barato de explicar:

1. **Depois de escrever um Draft** (steps 5 e 7): `GET` na versão, que precisa estar em `draft` e com o `version_id` esperado. Isso pega um Draft que outro processo (Console, outro CLI) colocou em build no meio do run. Nesse caso o CLI para com uma mensagem que diz isso, em vez de escrever por cima.
2. **Depois do build** (steps 6, 8 e 12): o estado final é `ready`, e o CLI grava `version_state` no azion.json.
3. **Depois da ativação** (step 15, `VerifyDeploy`):
   - a release está `ready` e com o papel de tráfego esperado;
   - `GET /deployments/{d}/resource_hierarchy` mostra a versão servida de cada recurso. Qualquer recurso da release (Application, Firewall, Custom Page ou dependência) ausente ou em versão diferente da que este run construiu é erro;
   - o Workload tem o binding para este Environment e este Deployment.

   Qualquer divergência vira erro com a diferença listada, e o deploy não é dado como concluído.

## 7. Falhas, interrupção e reexecução

| Onde falha | Efeito no tráfego | O que o CLI faz | No próximo run |
|---|---|---|---|
| Escrita de Draft (5, 7) | Nenhum | Para e grava o que já tinha feito. | Reaproveita os Drafts gravados. |
| Build de versão (6, 8) | Nenhum | Espera o nível inteiro e lista cada falha. | As versões em `error` são editáveis e são reaproveitadas. |
| Release build (12) | Nenhum | Para. Apaga a release em `error` que este run criou, para não gastar o limite de 20 Drafts de release. | Cria uma release nova. |
| Ativação (13) | Nenhum: na single_version a troca é atômica, e a release anterior continua ACTIVE. | Para com o erro da API. | Ativa de novo a mesma release Ready, sem rebuild. |
| Workload (14) ou verificação (15) | A release nova já está no ar. | Para e diz isso com todas as letras: o deploy está ativo, o que falhou foi o Workload ou a conferência. | Só repete o Workload e a verificação. |

**Interrupção (Ctrl-C, SIGTERM).** O deploy da v6 roda com `signal.NotifyContext`. Ao interromper:

- a espera para na hora;
- os builds que este run disparou e que ainda não terminaram recebem um cancel, por boa educação, com um contexto próprio de prazo curto (`context.WithoutCancel` mais `context.WithTimeout`);
- o azion.json é gravado;
- a mensagem diz se a interrupção foi antes ou depois da ativação.

**azion.json gravado a cada step**, como `WriteConfig` já faz na v4. Um processo morto no meio deixa o progresso registrado, e a seção 5 sabe continuar dali.

**Recursos de primeiro nível que saíram do manifest** (um Connector removido, por exemplo) não são apagados nesta versão. Eles podem estar em releases de outros Deployments. O CLI lista esses recursos como aviso. Apagar com segurança exige consultar `GET /resource_usage` antes, e isso é uma decisão separada (D4).

## 8. Contratos: azion.json e manifest da v6

Cada geração tem seu próprio tipo, como a v3 já tem `AzionApplicationOptionsV3`. A v6 não reaproveita o struct da v4: os campos mudam de significado, porque um id sozinho não diz qual versão está no ar.

**`contracts.AzionApplicationOptionsV6`** (azion.json):

- igual à v4: `name`, `bucket`, `preset`, `env`, `prefix`, `rotate-prefix`, `not-first-run`;
- `environment`: id e nome;
- `deployment`: id, nome e política;
- `release`: id, state e se foi ativada;
- `workload`: id, nome, domínios e url;
- por recurso (applications, functions, connectors, firewalls, custom pages, network lists, wafs): id, nome, `version_id` e `version_state`.

**`contracts.ManifestV6`** (gerado pelo bundler, P1):

- recursos, no formato dos requests do SDK da v6, como o `ManifestV4` faz com o da v4;
- `environment`: nome, opcional;
- `deployment`: política, opcional;
- `release.include`: dependências manuais, ou seja, Connectors e Network Lists que uma Function usa em tempo de execução e que a varredura da versão não enxerga;
- workloads com domínios.

O CLI valida o manifest no step 3, antes de qualquer chamada à API. O formato exato depende do bundler. Se ele chegar diferente, o ajuste é um adaptador em `ReadManifest`, sem mexer nos steps.

## 9. Layout no código

Segue a convenção registrada em `doc/plan.md`: a versão é uma subpasta de quem já é dono do código, sem árvore nova de primeiro nível.

| O quê | Onde |
|---|---|
| SDK | `go.mod`: `github.com/aziontech/azionapi-v6-go-sdk` |
| URL | `constants.ApiV6URL` e `viper.SetDefault("api_v6_url", ...)`, ao lado de `api_v4_url` |
| Clientes | `pkg/api/v6/<recurso>` (applications, connectors, custom_pages, firewalls, functions, network_lists, wafs, workloads, deployments, releases, environments), no mesmo formato de `pkg/api/applications/client.go` |
| Estados e erros | `pkg/api/v6/version` (`State`, `BuildError`, `Versioned`) |
| Espera | `pkg/poll` |
| Registry de clientes | `pkg/registry`: construtores v6 lendo `api_v6_url` |
| Applier do manifest | `pkg/manifest/v6` (`ResourceContextV6`, reconciliação de Drafts) |
| Pipeline e comando | `pkg/cmd/deploy_remote/v6` (`DeployStateV6`, `V6DeployPipeline`, dry-run, timer) e `pkg/cmd/deploy/v6` (o dispatcher, ver D1) |
| Mensagens | `messages/deploy-remote` e `messages/manifest`, em blocos com banner `v6`, como manda a regra de dono por versão |

**Registro do comando.** Hoje `cmdregistry.versionOf` serve a árvore da v4 inteira para contas v6. O deploy é o primeiro comando que diverge. Para isso:

- `cmdregistry` ganha `V6Only`;
- `V4Only` ganha uma variante exata (`V4Exact`), que não é servida à v6;
- a regra temporária passa de `versionOf` para `Entry.Supports`: uma entrada `V4Only` vale para a v6 até alguém declarar a versão v6 do mesmo comando;
- o teste de nomes únicos por versão passa a rodar também para a v6, e pega um `deploy` registrado duas vezes.

## 10. Superfície do comando

- **Flags:** as mesmas da v4, mais `--environment <nome>` (o nome `--env` já é usado para o arquivo `.env`), `--no-activate` (para depois do build da release) e `--timeout` (limite de cada espera; padrão 3m).
- **Progresso:** uma linha por transição de estado, com o tempo decorrido, por exemplo `application "site" (01JCA4M8): queued → building`. O spinner do runner continua; sob `--debug`, cada consulta vai para o log.
- **Saída final:** release id, Environment, Deployment, a versão de cada recurso e a URL (`*.azion.run`, ou a URL da release em versioned_urls).
- **`--format json`:** os mesmos dados, estruturados.
- **Resumo de timing:** um `TimingSummaryV6` próprio, indexado pelos nomes da seção 3, com o tempo de espera de build separado do tempo de chamada. Um teste prende os nomes dos steps aos campos do resumo, como `TestStepNamesAreRecordedByTheTimingSummary` faz na v4.
- **`--dry-run`:** derivado da lista de steps, como na v4. Diz quais recursos ganhariam Draft, quais builds e qual release, sem nenhuma requisição.

## 11. Testes

- **`pkg/poll`:** relógio falso via `Options.After`, sem `time.Sleep`. Casos: sucesso, timeout com último estado, `ctx` cancelado, 429 com e sem `Retry-After`, 5xx retentado, 4xx encerrando, estado desconhecido não encerrando.
- **Estados:** tabela com cada valor de `State` contra `Settled` e `Failed`.
- **Drafts (seção 5):** um caso por linha da tabela, com `httpmock`.
- **Builds de um nível:** um falha e os outros terminam; o erro junta todos; a ordem dos níveis é respeitada.
- **Release:** criar, build, ativar; build em `error` apaga a release do run; ativação falha e a anterior continua.
- **Interrupção:** `ctx` cancelado no meio da espera chama o cancel dos builds iniciados e grava o azion.json.
- **Verificação:** hierarquia diferente do esperado vira erro com a diferença.
- **Dry-run:** um `Transport` que reprova qualquer requisição e um `WriteAzionJsonContent` que reprova qualquer escrita, os mesmos guardas da v4.
- **Pipeline:** tabela de steps com ordem e `Applies`, e nomes presos ao resumo de timing.
- **Registry:** nomes únicos por versão incluindo a v6; contagem de comandos da v6.

## 12. Decisões

| # | Decisão | Situação |
|---|---|---|
| D1 | `azion deploy` (sem `--local`) na v6 | **Direção decidida:** a v6 terá deploy remoto próprio, com um script novo ou com o atual adaptado; qual dos dois ainda está em aberto. Até ele existir, a recomendação é que `azion deploy` em conta v6 rode o pipeline local e avise. |
| D2 | Esperar o build por padrão, ou devolver os ids e sair? | Em aberto. Recomendação: esperar por padrão; um `--no-wait` pode vir depois, se houver pedido. |
| D3 | Environment padrão | **Decidido:** Production. `--environment`, o azion.json ou o manifest podem apontar outro, nessa ordem de precedência. |
| D4 | Apagar recursos de primeiro nível que saíram do manifest | Em aberto. Recomendação: não nesta versão; avaliar depois com `GET /resource_usage`. |
| D5 | Intervalo de polling | **Decidido:** uma consulta a cada 5 segundos, sem backoff (era 1s; alterado em 05/11/2026). Limite de taxa informado: 200 requisições por segundo, a confirmar. Timeout padrão: 3 minutos por espera (era 50s), ajustável com `--timeout`. |
| D6 | De onde vem a URL de uma release em versioned_urls | Em aberto. O SDK não mostra o campo; perguntar ao time de API. |
| D7 | Versão de cada recurso na release | **Decidido:** o `version_id` construído no run (short_id como `A32XCE2E`). A API recusa `"latest"` com 400 (`version_id` precisa ser um short_id com 8+ caracteres). A release também nomeia todas as dependências do manifest, cada uma na versão construída, porque a API não as completa. A verificação final trata qualquer recurso da release ausente ou em outra versão como erro. |

As perguntas ao time de API da seção 13.2 de `doc/plan.md` continuam valendo. As que travam este plano:

- valores de `resource_type` na release (o step 11 precisa deles);
- `/v6` ou `/v4` (P2);
- se a release é finalizada por `.../build` ou pelo `POST` no próprio path.

## 13. Ordem de entrega

Cada linha é um PR que compila, tem teste e não muda nada para contas v3 e v4. Esforço a preencher pelo time.

| # | Entrega | Depende de | Esforço |
|---|---|---|---|
| 1 | SDK no `go.mod`, `api_v6_url`, `pkg/api/v6/version` e os clientes de `pkg/api/v6` com a interface `Versioned` | — | ? EW |
| 2 | `pkg/poll` | — | ? EW |
| 3 | `AzionApplicationOptionsV6` e `ManifestV6` | formato do bundler | ? EW |
| 4 | `pkg/manifest/v6`: Drafts e reconciliação (seção 5) | 1, 3 | ? EW |
| 5 | Steps de build com espera (6, 8) | 2, 4 | ? EW |
| 6 | Environment, Deployment, release, ativação, Workload e verificação (9 a 15) | 5, respostas de 13.2 | ? EW |
| 7 | Comando: `pkg/cmd/deploy_remote/v6`, dispatcher, `V6Only`/`V4Exact`, mensagens, timing, dry-run | 6 | ? EW |
| 8 | Deploy remoto da v6: script novo ou o atual adaptado | 7, D1 | ? EW |

## 14. Implementação (04/11/2026)

Entregas 1 a 7 da seção 13 implementadas, com testes. A 8 (deploy remoto) depende de D1. Nada foi testado contra a API real: os testes usam um fake da API v6 (`pkg/api/v6/v6test`) com respostas em sequência por rota.

| Peça | Onde |
|---|---|
| Clientes versionados, no mesmo formato dos da v4: `client.go`, um arquivo do recurso, `versions.go` e um arquivo por sub-recurso, cada chamada do SDK escrita por extenso com o tratamento de erro e de corpo da v4 (`LogAndRewindBodyV4`, `ErrorPerStatusCodeV4`) | `pkg/api/v6/{applications,connectors,custompages,firewalls,functions,networklists,wafs}` |
| Environments, Deployments (releases, ativação, hierarquia), Workloads | `pkg/api/v6/{environments,deployments,workloads}` |
| Estados, `BuildDone`, `StatusError`, `Retryable`, `BuildError`, `Versioned`, `Object`, `Item` | `pkg/api/v6/version` |
| Espera | `pkg/poll` |
| Contratos | `pkg/contracts/v6.go`, `utils/azionjson_v6.go` |
| Drafts, reconciliação, referências por nome, validação | `pkg/manifest/v6` |
| Pipeline, builds em dois níveis, release, Workload, verificação, timing, dry-run | `pkg/cmd/deploy_remote/v6` |
| Ponte para build, bucket e upload da v4 | `pkg/cmd/deploy_remote/shared_v6.go` |
| Comando | `pkg/cmd/deploy/v6`, registrado com `V6Only`; o da v4 passou a `V4Exact` |

**Onde a implementação difere do plano:**

- **azion.json:** `deployments` é uma lista com um Deployment por Environment, cada um com a sua release, em vez de um `deployment` e uma `release` únicos. Assim, um deploy em Staging não reaproveita o Deployment de Production. O campo `environment` virou o nome do Environment padrão do projeto.
- **Versão na release (D7):** `releaseResources` em `pkg/cmd/deploy_remote/v6/release.go` envia o `VersionID` de cada recurso construído no run, incluindo as dependências.
- **`resource_type`:** usa as chaves da hierarquia do Deployment (`application`, `firewall`, `custom_page`, `connector`, `network_list`, `waf`, `function`), no mapa `resourceTypes`. Também falta a confirmação do time de API.
- **Uma Application, um Firewall e uma Custom Page por manifest:** a hierarquia do Deployment tem um slot só para cada um. O `ReadManifest` recusa mais de um.
- **Workload:** o deploy aplica o primeiro Workload do manifest, como a v4.
- **Steps a mais:** `BundlerInit` e `ReadManifestAgain` (iguais aos da v4, porque o bundler precisa do nome do bucket no primeiro run) e `ReportOrphans` (D4: só avisa).
- **`--sync`:** não existe na v6, porque é sincronização de recursos da v4.
- **Erros:** iguais aos da v4. A única adição está nas leituras que a espera consulta (`GetVersion`, `VersionState`, `Deployments.Get`, `GetRelease`): o mesmo erro da v4 vai embrulhado em `version.StatusError`, que guarda o status HTTP e o `Retry-After`. A mensagem e o `errors.Is` com os sentinelas de `utils` não mudam. A espera retenta 429, 5xx e os sentinelas `ErrorInternalServerError` e `ErrorTimeoutAPICall`, que são o que a v4 devolve quando não há resposta.
- **`POST .../releases/{id}/build`:** o SDK exige um corpo; o cliente envia `{}`.

**Manifest do bundler (04/11/2026).** O `ManifestV6` segue agora a saída real do bundler com `--config-version 6` (pré-release `https://pkg.pr.new/@aziontech/bundler@666`; a saída completa de exemplo está em `pkg/manifest/v6/testdata/bundler666_full.json`). O que mudou:

- **Formato:** `version: 6` obrigatório; chaves `network_list`, `waf` e `firewall` (com `rules_engine`); regras da Application numa lista só com `phase`; `run_function` nomeia a Function, e o CLI usa a instância dela na mesma Application ou Firewall, como a v4; critérios `is_in_list`/`is_not_in_list` do Firewall nomeiam Network Lists; custom pages nomeiam o Connector. Tudo isso é resolvido para id antes de virar request do SDK.
- **Environments e Deployments vêm do manifest.** Cada Environment declarado é reaproveitado pelo nome (azion.json, nome exato, depois sem diferenciar maiúsculas) ou criado; a `deployment_policy` é imutável, e uma diferente para o deploy com erro. Cada Deployment é reaproveitado ou criado com o `binding_policy` do manifest e, sem `deployment_policy`, a do Environment a que está ligado. Há uma release por Deployment.
- **Bindings vêm do manifest** (`environment_name`, `deployment_name`, `domains`). Bindings de Environments que o run não toca são preservados.
- **`--environment` virou filtro:** só os Deployments ligados àquele Environment recebem release, e só os bindings dele mudam. No primeiro run ele também substitui `$ENVIRONMENT_NAME` no azion.config (padrão: Production, D3).
- **Manifest sem Environments e Deployments:** o deploy usa Production (que precisa existir) e um Deployment FLEXIBLE com o nome do projeto.
- **`--no-activate` não mexe no Workload:** ligar um Deployment sem release ativa deixaria o Environment sem nada no ar.
- **Bundler da v6:** `vulcan.NewVulcanV6` roda o pacote em `BundlerV6Package` (hoje o pré-release; `AZIONCLI_BUNDLER_V6` sobrescreve) com `--config-version 6`. O build chamado pelo deploy da v6 não lê nem grava mais o azion.json no formato da v4, o que antes sobrescreveria o arquivo da v6.
- **Removido:** `release.include` (inventado por mim; o bundler não tem). Fica a pergunta ao time da lib: como o manifest declara Connectors e Network Lists que uma Function usa em tempo de execução, que a release precisa incluir à mão?

**Antes de liberar:**

- confirmar com o time de API `resource_type`, `/v6 contra `/v4` (P2) e como a release é finalizada;
- trocar `BundlerV6Package` pela versão publicada do bundler quando ela sair;
- secret `AZION_API_V6_STAGE` no GitHub para o build de stage;
- um deploy de ponta a ponta em uma conta v6 de stage.

