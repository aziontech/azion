Azion CLI — Arquitetura multi-versão e plano de preparação para a API v6
Status: proposta — para decisão
Data: 10 set 2026
Autoria: Patrick Selau Menoti
Revisão: Jose Filho, Magnun A V F, Pablo Diehl
Quadro-resumo
O problema em três números: 47% da árvore de comandos de legado é redundante (5.966 LOC), 17% das mudanças no CLI precisam ser feitas duas vezes, e um terceiro fork multiplicaria esse custo por três. Os três são medidos no repositório; os comandos de verificação estão no apêndice.
Fases de execução

Fase
O que entrega
Tier
Esforço*
Depende de
Precisa do schema da v6?
0
Fundações:
Fim do estado global de flags, resolução de versão de conta com cache
Discovery da v6
—
3 turnos
—
Não
1
Consolidação da duplicação redundante
Tier 1
4 turnos
Fase 0
Não
2
Isolamento dos recursos exclusivos de cada geração
Tier 2
8 turnos
Fase 1
Não
3
Extração da pipeline de deploy para um runner único
Tier 3
6 turnos
Fase 2
Não


Total da melhoria (fases 0 a 3)


21 turnos


Não
4
API v6: novo client, recursos e pipeline. Nenhuma árvore de comandos nova (se seguir algo similar ao que já existe)
todos
Não é possível estimar ainda.
Fases 0 a 3 + respostas A1 e A2
Sim


*Estimativas de esforço e prazo são estimativas abstratas levando em consideração um escopo em aberto a ser executado por um único desenvolvedor. Detalhamento completo na seção 12.2
Os três tiers
Tier
O que é
Escopo medido
Abstração
Tier 1
Comandos agnósticos de versão: uma única cópia, e o cliente de API vem de um registry
59 arquivos, 5.966 LOC redundantes hoje
Interface fina, que apenas documenta um contrato que os packages já cumprem
Tier 2
Recursos que existem em uma só geração, registrados por versão
32 arquivos exclusivos de v3 (3.161 LOC), mais o conjunto exclusivo de v4
Nenhuma: conceitos diferentes permanecem diferentes
Tier 3
Pipeline de deploy versionada, com a lista de steps fornecida por cada geração
15 arquivos, 3.518 LOC
A única interface real do plano


Detalhamento na seção 10. O motivo de não haver uma única abstração para tudo está na seção 5.
Decisões solicitadas: aprovar o sequenciamento (melhorias antes da v6), confirmar a alocação de 1 ou 2 engenheiros, endossar as 8 solicitações formais ao time de API, e confirmar que clientes de v3 e v4 não são impactados. Detalhamento nas seções 8 e 13.


Parte I — Resumo
1. A situação atual
O Azion CLI dá suporte a duas gerações da nossa API ao mesmo tempo. E faz isso mantendo duas cópias completas de si mesmo dentro de um único binário:
Camada
API v4 (atual)
API v3 (legado)
Comandos
pkg/cmd/
pkg/v3commands/
Clientes de API
pkg/api/
pkg/v3api/
Manifest de deploy
pkg/manifest/
pkg/v3manifest/
SDK gerado
azionapi-v4-go-sdk
azionapi-go-sdk


Na inicialização, o CLI consulta nosso serviço de SSO para saber a qual geração pertence a conta autenticada e, então, monta uma das duas árvores de comandos. O fork foi introduzido em 26/05/2025 (fc1639cf "chore: add retrocompatibility") como um atalho deliberado e pragmático para lançar a v4 sem quebrar os clientes existentes. Funcionou, e a seção 2 expõe as restrições que o tornaram a escolha correta na época. Também criou um custo que pagamos a cada mudança.
A API v6 está a caminho e vai introduzir um novo modelo de deploy e de criação de recursos. O caminho natural, copiar a árvore uma terceira vez, levaria o CLI de duas implementações paralelas para três. Este documento argumenta contra isso, quantifica o motivo e propõe uma alternativa.
2. Por que a arquitetura atual foi feita assim
O fork não foi um descuido. Foi a única saída viável diante de algumas restrições que apareceram em uma ordem não ideal.
Retrocompatibilidade não estava no plano original da v4. Ela chegou como requisito depois, com o trabalho da v4 já em andamento. O histórico do repositório registra isso: o SDK da v4 entrou no go.mod em 19/03/2025 e os primeiros comandos v4 foram entregues em abril (feat: add edge-application v4 commands, 23/04/2025). O commit do fork, fc1639cf "chore: add retrocompatibility", é de 26/05/2025, 68 dias depois do início da v4.
E o schema da API v4 ainda estava em movimento. Este é o fator decisivo, e é mensurável: o go.mod fixou diversas versões distintas do SDK ao longo do projeto, de v0.22.0 a v0.266.0. Recursos centrais do modelo v4 continuaram mudando bem depois do lançamento; pkg/api/applications e pkg/api/connector só foram criados em 24/09/2025. Não se projeta uma abstração compartilhada estável contra um schema que ainda está se definindo: a abstração erra o alvo, e refazê-la custa mais do que duplicar.
Diante dessas restrições, duplicar a árvore foi a escolha correta. Ela entregou o que precisava entregar: a v4 saiu na data e nenhum cliente de v3 parou de funcionar.
O que é diferente agora. Três coisas mudaram, e é por isso que este documento propõe fazer o trabalho de outra forma:
Sabemos da v6 antes de ela existir. Em 2025, a retrocompatibilidade virou requisito no meio da execução. Agora estamos planejando com antecedência, sem uma data de lançamento pressionando o design.
Temos a experiência de duas gerações para aprender. As quatro áreas de divergência real - vocabulário de recursos, schema do manifest, orquestração de deploy e cliente de SDK - só ficaram visíveis depois de conviver com o fork. São elas que fundamentam o design em três tiers.
E, principalmente: o trabalho de arquitetura não depende do schema da v6. As fases 0 a 3 operam inteiramente sobre código v3/v4 que já existe. Nenhuma delas precisa conhecer uma linha do schema da v6, só a fase 4 precisa. Ou seja, a restrição que impediu o design correto em 2025 não nos impede agora, mesmo com a v6 ainda indefinida.
3. Quanto o fork nos custa hoje
Os números abaixo foram medidos no repositório. Não são estimativas.
Quase metade da árvore de legado é redundante. Das 12.645 linhas em pkg/v3commands/:
Categoria
Arquivos
Linhas
Proporção
Idênticos byte a byte ou com diferença apenas na linha de import
40
3.901
31%
Pequenas diferenças mecânicas (16 a 100 linhas)
19
2.065
16%
Comportamento genuinamente diferente
15
3.518
28%
Recursos que existem apenas na v3
32
3.161
25%


47% da árvore de comandos de legado, 5.966 linhas, é código redundante que existe puramente por causa de como o fork foi construído. Não carrega nenhuma lógica específica de v3. Parte disso vai além de redundante: pkg/v3api/variables, pkg/v3api/personal_token e pkg/v3api/graphql são idênticos byte a byte aos seus equivalentes de v4, 477 linhas de duplicação exata, sem nenhuma diferença.
17% de todas as mudanças no CLI precisam ser feitas duas vezes. Nos últimos 12 meses, 35 de 206 commits tiveram que alterar as duas árvores para entregar uma única mudança. Cada uma é aplicada à mão, duas vezes, sem nada que verifique se as duas edições coincidem.
O custo não é paridade de features, é infraestrutura compartilhada. A árvore de v3 está deliberadamente sob uma política de somente correções de bugs: novas features da API não são retroportadas. As diferenças de flags no deploy (a v3 não tem --workers, --writable-bucket nem --skip-framework-build) são essa política funcionando como projetado, não divergência acidental.
A composição desses commits espelhados mostra tanto que a política se sustenta quanto onde está o custo real. Apenas 2 dos 35 são fix; 25 são chore, refactor ou tests. Os 4 commits feat espelhados são todos de ferramental compartilhado, não de recursos da API: --alias-env (uma flag de build do Bundler), profiles (gestão de credenciais), rotate-prefix (templates de deploy) e uma opção de workflow no link.
Este é o ponto estrutural. A divergência de features é governada, e a governança funciona. O que não pode ser congelado junto com a superfície de API da v3 é a infraestrutura que todas as versões compartilham: build, a mecânica de deploy, credenciais, output, testes, atualizações de dependência. Essa camada é portada à mão a cada vez, e é exatamente a camada que a refatoração em três tiers consolida.
Um exemplo concreto. O PR #1539, feat: add alias-env flag, adicionou uma única flag booleana. Ele alterou 20 arquivos: 10 na árvore de v4, 5 na árvore de v3 e 5 arquivos de mensagens. O arquivo pkg/v3commands/build/run_test.go nesse PR é uma cópia literal de 57 linhas do teste de v4. Uma flag, duas implementações, duas suítes de teste.
4. Quanto custaria um terceiro fork
Copiar a árvore novamente para a v6 iria:
Adicionar entre 15.000 e 25.000 linhas de código majoritariamente duplicado a uma base de 96.569 linhas.
Transformar a taxa de 17% de manutenção duplicada em um custo triplo. O padrão do PR #1539 passa a ser algo entre 28 e 30 arquivos para uma única flag.
Triplicar o trabalho manual de portar infraestrutura compartilhada. A política de somente correções de bugs contém a divergência de features; ela não faz nada pela camada de housekeeping, que são 25 dos 35 commits espelhados. Com três árvores, toda mudança de build, credencial, output e dependência é aplicada três vezes à mão.
Tornar a consolidação futura muito mais difícil. Reconciliar três árvores divergentes é significativamente pior do que duas, e a janela para fazer isso de forma barata se fecha no momento em que a v6 for lançada.
5. Recomendação
Primeiro melhorar o código existente da arquitetura de duas versões e, depois, adicionar a v6 na estrutura resultante. Não fazer um terceiro fork.
Propomos uma arquitetura graduada em três tiers, em vez de uma única grande abstração. O princípio de design é deliberadamente conservador:
Compartilhar código apenas onde ele já é idêntico. Isolar o código que genuinamente difere, sem forçá-lo por uma interface comum. Construir uma abstração real para exatamente uma coisa, a pipeline de deploy, porque é o único lugar em que as versões diferem em comportamento, e não em vocabulário.
Isso importa porque uma abstração completa é tentadora e errada. domain da v3 e workload da v4 não são o mesmo conceito; origin e connector também não. Forçar ambos por uma única interface produz o modo de falha clássico: métodos não implementados, erros de "não suportado nesta versão" e verificações de versão voltando para dentro da abstração que deveria eliminá-las. A abordagem graduada recusa essa troca.
Restrição não negociável: clientes com setups funcionais de v3 e v4 devem continuar funcionando exatamente como hoje, durante e depois da refatoração. As fases de refatoração preservam o comportamento de forma estrita. Melhorias que mudariam comportamento observável ficam para o próximo major release do CLI, onde podem ser comunicadas como tal.
6. Investimento e cronograma
Considera 1 ou 2 engenheiros dedicados, conforme acordado.
Fase
Escopo
Esforço
0
Fundações + discovery de v6
3
1
Eliminar a duplicação redundante (Tier 1)
4
2
Isolar recursos exclusivos de cada versão (Tier 2)
8
3
Extrair a pipeline de deploy (Tier 3)
6


Total da refatoração
21


Tempo de calendário:
1 engenheiro dedicado: 2 semanas + 1 turno
O que isso compra na entrega da v6:
Abordagem
Custo de implementação da v6
Custo contínuo
Terceiro fork (sem refatoração)
?
Custo permanente de manutenção em 3 vias
Dentro da estrutura refatorada
?
Uma camada compartilhada, um plug-in por versão


A refatoração se paga em cerca de metade já na entrega da v6, e o restante nas duas ou três primeiras mudanças transversais depois disso.
As fases 1, 2 e 3 podem ser entregues de forma independente. Cada uma termina com um CLI funcional e liberado. Se as prioridades mudarem, o trabalho pode ser interrompido no limite de uma fase sem deixar a base de código em estado pior do que estava no início.
7. Riscos
Risco
Severidade
Mitigação
A mudança causar regressão em um setup de cliente que funciona
Alta
Toda fase preserva o comportamento de forma estrita; as fases são entregues de forma independente e podem ser revertidas individualmente. A fase 0 entrega a rede de segurança antes de qualquer consolidação.
O schema real da v6 invalidar o design
Média
A fase 0 é uma fase de discovery que produz requisitos formais para o time de API. A abstração é deliberadamente fina e restrita à pipeline de deploy; a camada com maior chance de precisar de retrabalho é também a menor.
A melhoria competir com a entrega de features
Média
Decisão explícita de alocação solicitada abaixo. Os limites entre fases permitem intercalar com trabalho de features a um custo conhecido.
O cronograma da v6 encurtar e forçar um fork de qualquer forma
Média
Só a fase 1 já remove a maior parte da redundância e torna até um fork de v6 mais barato. Sequenciar a fase 1 primeiro preserva essa opção.
Scope creep virando um rewrite completo
Média
O design graduado é a própria salvaguarda: três tiers explícitos, com critérios escritos para decidir a qual tier cada trecho de código pertence. O que não se encaixa em nenhum tier está fora de escopo.


O risco de não agir é triplicarmos permanentemente o custo de toda mudança na camada de infraestrutura compartilhada: a camada que concentra a maior parte do nosso trabalho no CLI, e a única que nenhuma política de versionamento consegue congelar.
8. Decisões solicitadas
Aprovar o sequenciamento "melhoria primeiro" — a implementação da v6 começa após a conclusão da fase 1, não antes.
Confirmar a alocação — 1 ou 2 engenheiros dedicados, o que define o calendário de [tempo] de trabalho.
Endossar as solicitações da fase 0 ao time de API (seção 13) — o CLI precisa de compromissos específicos sobre a v6 antes de poder construir sobre ela.
Confirmar a restrição de compatibilidade — clientes atuais de v3 e v4 não são impactados; mudanças de comportamento ficam para o próximo major release do CLI.


Parte II — Arquitetura Técnica
9. Estado atual, medido
9.1 Escala do repositório
Todas as contagens desta seção são em LOC (lines of code, linhas de código), medidas com wc -l sobre arquivos .go, conforme os comandos do apêndice.
Métrica
Valor
Total de LOC em Go
96.569
LOC de testes
34.862
Árvore de v4 (pkg/cmd + pkg/api + pkg/manifest), sem testes
36.831
Árvore de v3 (pkg/v3commands + pkg/v3api + pkg/v3manifest), sem testes
15.947


9.2 O que já é compartilhado
A base de código já demonstrou que uma camada compartilhada funciona. As duas árvores consomem:
messages/ · pkg/contracts · pkg/output · pkg/cmdutil · pkg/token · pkg/config · pkg/iostreams · pkg/logger · pkg/vulcan · utils/
Mais revelador ainda: pkg/cmd/rollback e pkg/cmd/profiles são registrados na árvore de comandos de v3 diretamente a partir do package de v4: veja v3rollback em root.go:49 e a chamada de profiles.NewCmd dentro de setV3Cmds. Comandos agnósticos de versão já existem. A proposta generaliza um padrão estabelecido; não introduz um novo.
O package compartilhado messages/ mostra onde o design atual para no meio do caminho. Os dois arquivos deploy.go importam o mesmo messages/deploy, e esse único arquivo guarda strings que as duas árvores usam (AliasEnvFlag) ao lado de strings que só a v4 registra (WorkersFlag, WritableBucketFlag, SkipFrameworkBuild), sem nada marcando quais são quais. Nada quebra hoje, porque uma flag não registrada simplesmente não renderiza texto de ajuda. Mas significa que a pessoa que desenvolve não tem nenhum sinal estrutural sobre qual versão é dona de cada string. O compartilhamento parou nas strings e nunca chegou ao código que as usa. Um caso concreto: os dois comandos de application importam esse mesmo messages/create/application, que define Usage = "application". A v4 usa esse valor (Use: msg.Usage); a v3 o ignora e escreve Use: "edge-application" hardcoded em edge_application.go:66. Uma das árvores importa um package de strings e não usa o valor principal dele.
9.3 Duplicação classificada
Os 106 arquivos sem testes em pkg/v3commands/, classificados por quanto divergiram do seu equivalente em pkg/cmd/:
Faixa
Arquivos
LOC
Exemplos
Tier alvo
≤15 linhas de diff (apenas o caminho do import)
40
3.901
CRUD de origin, personal_token, variables, domain, além de logs/cells, logs/http, login, version, completion, whoami
Tier 1
16 a 100 linhas de diff
19
2.065
init/*, link/*, purge/*, dispatchers create.go/list.go/update.go
Tier 1 após parametrização
>100 linhas de diff
15
3.518
deploy_remote/requests.go (462), sync/tasks.go (387), deploy/deploy.go (138), deploy/upload.go (199), `create
update/cache_setting` (~200 cada)
Sem equivalente em v4
32
3.161
domain, edge_applications, edge_function, edge_storage
Tier 2


Oito arquivos são idênticos byte a byte: build/utils.go, completion/completion.go, deploy/requests.go, deploy/scriptrunner.go, login/login_mock.go, logs/logs.go, unlink/utils.go, version/version.go.
Na camada de clientes de API, a redundância é ainda mais evidente. Estes packages de pkg/v3api são 100% idênticos byte a byte aos seus equivalentes em pkg/api:
Package
Arquivos
LOC
v3api/variables
3
160
v3api/personal_token
5
148
v3api/graphql (+ cells, http)
3
169


11
477


477 linhas removíveis apenas trocando caminhos de import. Sem interface, sem abstração, sem trabalho de design; esta é a primeira evidência prática da fase 1.
A descoberta estrutural crítica: a divergência genuína está confinada a quatro preocupações: o vocabulário de recursos (domain/origin versus workload/connector), o schema do manifest, a orquestração de deploy e o cliente de SDK. Todo o resto é acidental.
9.4 Defeitos encontrados nas camadas que vamos redesenhar
Os itens abaixo são bugs pré-existentes descobertos durante a análise de três camadas: a seleção de versão (D1 a D3), o estado das flags (D4) e a pipeline de deploy (D5). Eles importam porque são exatamente as camadas que vamos redesenhar: cada defeito é corrigido como subproduto de fazer isso corretamente, e cada um seria herdado pela v6 caso contrário.
D1 — --token não afeta a seleção de versão. cmd/azion/main.go:21-30 lê o token armazenado do profile ativo e popula o Viper antes de o Cobra existir. CmdRoot() então consulta a conta e monta a árvore de comandos usando esse token, em root.go:222. A flag --token é registrada em root.go:121, mas só é consumida em persistentPreRunE → doPreCommandCheck, que o Cobra executa depois de a árvore já estar montada.
Consequência: azion --token <token-de-uma-conta-v4> list workload seleciona a árvore de comandos a partir de qualquer token que esteja em disco. Para pipelines de CI que se autenticam exclusivamente via --token, o CLI pode apresentar os comandos da geração errada por completo. A flag --config tem o mesmo defeito de ordenação.
D2 — Toda invocação faz uma chamada de rede bloqueante, sem cache. HasBlockAPIV4Flag (utils.go:20) emite um GET síncrono para sso.azion.com/api/account/info dentro de CmdRoot(), antes de qualquer parsing de argumentos. Tanto azion --version quanto azion --help pagam esse custo. O resultado nunca é persistido.
D3 — Falha na consulta faz downgrade silencioso para v3. root.go:214-224: qualquer erro na consulta — timeout, falha de DNS, 500, token expirado — cai em fact.apiVersion = "v3" e setV3Cmds. Um cliente de v4 em uma conexão instável recebe a árvore de comandos de v3 apontada para endpoints de v3. A falha é registrada apenas em nível de debug.
D4 — O estado das flags são variáveis globais mutáveis de package. pkg/cmd/deploy/deploy.go:59-74 declara Path, Auto, NoPrompt, SkipBuild, Sync, DryRun, Result e outras como variáveis de package, espelhadas na árvore de v3. Um comando não pode ser instanciado duas vezes no mesmo processo, o que bloqueia tanto a consolidação do Tier 1 quanto qualquer teste table-driven. Este é um pré-requisito rígido, não uma limpeza opcional.
D5 — A pipeline do manifest applier está escrita manualmente, de forma repetida. pkg/manifest/manifest.go:80-210 invoca ApplyFunctions, ApplyFunctionInstances, ApplyEdgeApplication, ApplyCacheSettings, ApplyConnectors, ApplyRulesEngine, ApplyWorkloads, ApplyWorkloadDeployments, ApplyFirewalls e ApplyPurge como dez blocos quase idênticos, cada um repetindo o mesmo boilerplate de verificação de tamanho, time.Now(), verificação de erro e GlobalTimingCallback. A ordenação, que é conhecimento de domínio real sobre dependências entre recursos, está implícita na ordem do código-fonte de uma função imperativa.
Esta é também a boa notícia: resource_applier.go já tem o formato correto. ResourceContext guarda os clientes e os mapas de IDs, e expõe um conjunto de steps Apply*. O ponto de extensão de que a v6 precisa já existe; só não está declarado como tal.
10. Arquitetura-alvo: três tiers
Tier 1 — Comandos agnósticos de versão (uma cópia)
Critério: as implementações de v3 e v4 diferem apenas em qual package de cliente de API importam, ou por substituições mecânicas.
Escopo: os 40 arquivos quase duplicados (3.901 LOC) mais os 19 arquivos da faixa intermediária (2.065 LOC) — variables, personal_token, origin, logs, login, logout, version, completion, whoami, reset, unlink e os dispatchers de create/list/update/delete/describe.
Mecanismo. Um enum de versão na factory e um registry de clientes:
// pkg/apiversion/apiversion.go
type Version string
const (
    V3 Version = "v3"
    V4 Version = "v4"
    V6 Version = "v6"   // adicionado na fase 4
)
 
// pkg/cmdutil/factory.go
type Factory struct {
    HttpClient *http.Client
    IOStreams  *iostreams.IOStreams
    Config     config.Config
    APIVersion apiversion.Version   // resolvido uma vez, na inicialização
    Flags
}
Cada comando compartilhado declara a interface restrita de que precisa e a obtém do registry:
// pkg/cmd/create/variables/variables.go   — UMA cópia, substitui duas
type VariablesCreator interface {
    Create(ctx context.Context, req api.CreateRequest) (api.VariableResponse, error)
}
 
func NewCmd(f *cmdutil.Factory) *cobra.Command {
    return newCobraCmd(f, registry.Variables(f))
}
 
// pkg/registry/variables.go
func Variables(f *cmdutil.Factory) VariablesCreator {
    switch f.APIVersion {
    case apiversion.V3: return v3api.NewClient(f.HttpClient, f.Config.GetString("api_url"), tok(f))
    case apiversion.V4: return api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), tok(f))
    }
}
Por que isso é barato: para os 40 arquivos quase duplicados, os dois packages de cliente de API já expõem assinaturas de método idênticas — é precisamente por isso que o diff é de uma linha de import. Para os 11 arquivos idênticos byte a byte, nem interface é necessária; o registry devolve o mesmo tipo para as duas versões. As interfaces aqui documentam um contrato que já existe, não introduzem um design novo.
Isso também substitui os 181 call sites escritos à mão com api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token")) — a string "api_v4_url" hardcoded já é, por si só, um problema latente para a v6.
Tier 2 — Recursos por versão (sem tentativa de abstração)
Critério: o recurso existe em uma geração da API e não nas outras, ou seu formato difere o suficiente para que uma interface compartilhada vazasse.
Escopo:
Exclusivos de v3 (32 arquivos, 3.161 LOC): domain, edge_applications, edge_function, edge_storage
Exclusivos de v4: workload, workload_deployment, connector, waf, waf_exceptions, data_stream, firewall*, crl, csr, custom_pages, device_groups, dns*, kv, network_list
Exclusivos de v6: o que a v6 vier a introduzir
pkg/resources/
  v3/   domain/  edge_applications/  edge_function/  edge_storage/
  v4/   workloads/  connector/  waf/  data_stream/  firewall/  ...
  v6/   <definido na fase 4>
Mecanismo. A divisão por versão não está na raiz. Medindo contra as duas listas de AddCommand em root.go:134-186, 22 dos 25 comandos de topo existem nas duas gerações; apenas clone, warmup e config são exclusivos de v4. A divisão real está um nível abaixo, dentro dos dispatchers:
azion create
filhos
v4 (create.go)
27
v3 (create.go)
10
mesmo nome de CLI nas duas
6 (rules-engine, personal-token, origin, cache-setting, variables, profile)


Uma lista plana na raiz não consegue expressar isso, porque os filhos precisam ser pendurados no pai create. Então o registry compõe em dois níveis. Uma tabela declara, por resource, em quais versões ele existe:
// pkg/registry/table.go
type entry struct {
    versions []apiversion.Version
    new      func(*cmdutil.Factory) *cobra.Command
}
 
var createEntries = []entry{
    // Tier 1 — uma implementação, todas as versões
    {all,      variables.NewCmd},
    {all,      personalToken.NewCmd},
    {all,      origin.NewCmd},
    {all,      cacheSetting.NewCmd},
    {all,      profile.NewCmd},      // já compartilhado pelas duas árvores hoje
 
    // Tier 2 — exclusivos de uma geração, sem tentativa de abstração
    {only(V3), v3domain.NewCmd},
    {only(V3), v3edgeApps.NewCmd},
    {only(V4), applications.NewCmd},
    {only(V4), workloads.NewCmd},
    {only(V4), connector.NewCmd},
    // ...
}
Cada dispatcher passa então a ser uma casca única que se preenche a partir da tabela:
// pkg/cmd/create/create.go — UMA cópia, substitui duas
func NewCmd(f *cmdutil.Factory) *cobra.Command {
    cmd := &cobra.Command{
        Use:   msg.Usage,
        Short: msg.ShortDescription,
        Long:  msg.LongDescription,
        RunE:  func(c *cobra.Command, _ []string) error { return c.Help() },
    }
    for _, child := range registry.Children("create", f) {   // filtrado por f.APIVersion
        cmd.AddCommand(child)
    }
    cmd.Example = examplesFrom(cmd)   // gerado a partir dos filhos, não hardcoded
    cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
    return cmd
}
Dois detalhes que vale explicitar. A versão não é parâmetro: depois da fase 0 ela vive em f.APIVersion, então o registry a lê em vez de recebê-la. E o Example precisa ser gerado — hoje ele é a única razão pela qual as duas cascas de create diferem, já que o heredoc da v3 lista edge-application onde o da v4 lista application. Derivá-lo dos filhos registrados torna as cascas genuinamente idênticas e remove uma fonte de divergência.
Explicitamente não feito aqui: nenhuma tentativa de unificar domain com workload, ou origin com connector. São conceitos diferentes e são modelados como tal. Esta é a concessão deliberada que mantém a abstração honesta, e é por isso que este é um design de três tiers, e não uma única interface de provider.
Tier 3 — Pipeline de deploy versionada (a única abstração real)
Critério: as versões diferem em comportamento e ordenação, não apenas em tipos. São os 15 arquivos divergentes (3.518 LOC) — deploy, deploy_remote, sync, manifest.
Mecanismo. Substituir a sequência escrita manualmente em manifest.go:80-210 por uma lista declarada de steps e um runner compartilhado:
// pkg/pipeline/pipeline.go
type Step struct {
    Name    string                        // usado para timing + logs + dry-run
    Applies func(*Plan) bool              // substitui as verificações inline de len(...) > 0
    Run     func(context.Context, *Plan) error
}
 
type Pipeline interface {
    Name()  string
    Steps() []Step
}
 
// O runner passa a ser dono de tudo o que hoje é copiado e colado dez vezes:
// spinner, timing callbacks, logging estruturado, renderização de dry-run,
// wrapping de erros e limpeza de recursos órfãos.
func Run(ctx context.Context, p Pipeline, plan *Plan) error
Cada versão fornece sua própria lista de steps:
// pkg/pipeline/v4/pipeline.go
func (p *V4Pipeline) Steps() []Step {
    return []Step{
        {Name: "functions",            Applies: hasFunctions,   Run: p.applyFunctions},
        {Name: "function-instances",   Applies: hasInstances,    Run: p.applyFunctionInstances},
        {Name: "application",          Applies: hasApplication,  Run: p.applyApplication},
        {Name: "cache-settings",       Applies: hasCache,        Run: p.applyCacheSettings},
        {Name: "connectors",           Applies: hasConnectors,   Run: p.applyConnectors},
        {Name: "rules-engine",         Applies: hasRules,        Run: p.applyRules},
        {Name: "workloads",            Applies: hasWorkloads,    Run: p.applyWorkloads},
        {Name: "workload-deployments", Applies: hasDeployments,  Run: p.applyDeployments},
        {Name: "firewalls",            Applies: hasFirewalls,    Run: p.applyFirewalls},
        {Name: "purge",                Applies: hasPurge,        Run: p.applyPurge},
    }
}
Por que esta é a abstração correta e as outras não são. A lista de steps torna a ordenação de dependências entre recursos — conhecimento de domínio genuíno, hoje implícito na ordem do código-fonte — explícita e revisável. Dá ao --dry-run uma implementação real de graça. Dá timing por step sem dez cópias das mesmas seis linhas. E, criticamente: a v6 mudar o modelo de deploy passa a significar fornecer uma lista de steps diferente, não fazer um fork do CLI. Este é o único lugar onde a "nova ideia de deploy e criação de recursos" aterriza.
A interface é intencionalmente mínima. Step tem um nome, um predicado e uma função. Sem taxonomia de recursos, sem vocabulário comum, sem schema compartilhado. Uma pipeline de v6 cujos steps não tenham nada em comum com os da v4 continua se encaixando.
11. Redesenho da resolução de versão
Hoje: uma chamada de rede por invocação, sem cache e incondicional; um resultado booleano; downgrade silencioso em caso de falha (defeitos D1 a D3).
11.1 Ordem de resolução
Quando a v6 for lançada, a precedência passa a ser v6 → v4 → v3: resolver para a geração mais recente à qual a conta tem direito. A implementação atual devolve um bool a partir de uma verificação de duas flags; ela passa a ser uma avaliação ordenada que devolve um enum, de modo que adicionar uma geração é um novo case, e não um novo ramo de lógica booleana.
11.2 Cache da versão resolvida no profile
Persistir a versão resolvida no arquivo de settings do profile, eliminando a chamada ao SSO em cada invocação (D2).
Já existe um precedente exatamente para isso na base de código. pkg/token/models.go:32-44 tem LastCheck time.Time e LastVulcanVersion, e pre_command.go:211 implementa um TTL de 24 horas sobre esse campo:
if time.Since(settings.LastCheck) < 24*time.Hour && !settings.LastCheck.IsZero() {
O cache de versão reutiliza esse formato:
type Settings struct {
    // ... campos existentes ...
    APIVersion          string    // "v3" | "v4" | "v6"
    APIVersionCheckedAt time.Time
    APIVersionTokenHash string    // vincula o cache à credencial que o produziu
}
Mecanismos de segurança — o cache nunca deve deixar na mão um cliente cuja conta foi promovida ou rebaixada no servidor:
Mecanismo
Disparo
Comportamento
TTL
Cache com mais de 24h
Nova consulta; reutiliza o padrão de pre_command.go:211
Vínculo com a credencial
APIVersionTokenHash ≠ hash do token em uso
Invalida e consulta de novo. Também corrige o D1: um --token diferente do armazenado força a resolução contra a conta correta
Invalidação por ciclo de vida
login, logout, troca de profile, mudança de --config
Limpa incondicionalmente
Correção pelo servidor
Qualquer resposta da API indicando que a versão não é permitida para a conta (403/409 ou o sinal equivalente na v6)
Invalida, consulta de novo e informa ao usuário que a geração da conta mudou e que ele deve executar novamente. Requer um sinal de erro distinguível — veja a solicitação A4
Override explícito
Flag --api-version / configuração de profile / AZIONCLI_API_VERSION
Ignora o cache e a consulta por completo. Torna o CI determinístico e o resolver testável
Reset manual
Subcomando azion config (ou azion login --refresh)
Limpa a versão em cache


Restrição de implementação (a partir do D1). A árvore de comandos é montada antes de o Cobra parsear qualquer coisa, então --api-version, --token e --config não podem ser lidos pelo binding normal de flags no momento da resolução. A resolução precisa de uma passada de pré-parsing sobre os.Args para essas três flags, antes de CmdRoot() montar a árvore. É uma mudança pequena e contida, e é a correção do D1, que hoje é um bug ativo para CI autenticado via --token.
Sobre o downgrade silencioso (D3). Falhar de forma explícita, em vez de selecionar v3 silenciosamente, é o comportamento correto, mas é uma mudança de comportamento observável. Sequenciamento:
Durante a refatoração: o cache e o override entram; o fallback é restringido — passa a valer apenas quando não há versão em cache, e emite um aviso visível em vez de um log de debug.
No próximo major release do CLI: falha na consulta sem versão em cache passa a ser um erro fatal com mensagem acionável, comunicado como breaking change.
Isso respeita a restrição de que clientes funcionais de v3 e v4 não sejam impactados.
12. Plano de fases
Cada fase termina com um CLI funcional e liberado. As fases podem ser revertidas individualmente.
Fase 0 — fundações e discovery de v6 · 3 turnos
Duas trilhas paralelas.
Trilha A — discovery de v6 (seção 13). Produzir um documento formal de requisitos para o time de API e obter respostas por escrito. É uma entrada obrigatória da fase 4 e deve começar imediatamente; é majoritariamente tempo de espera, e é por isso que roda primeiro e em paralelo.
Trilha B — pré-requisitos. Trabalho de base que preserva comportamento:
Eliminar as variáveis globais de flag em nível de package (D4). Mover Path, Auto, Sync, Result e as demais para structs de comando nas duas árvores. Desbloqueia todo o resto.
Introduzir apiversion.Version e Factory.APIVersion. Substituir o booleano por um enum; setCmds/setV3Cmds continuam consumindo-o. Sem mudança de comportamento.
Passada de pré-parsing para --api-version/--token/--config. Corrige o D1.
Cache de versão com TTL + vínculo com a credencial. Remove a chamada ao SSO por invocação (D2). Restringe o fallback silencioso e o torna visível (D3 parcial).
Critérios de saída: nenhum estado de flag em nível de package; a resolução de versão é baseada em enum, tem cache, aceita override e está correta em relação a --token; árvores de comandos inalteradas.
Fase 1 — consolidação do Tier 1 · 4 turnos
Apagar os packages de v3api idênticos byte a byte (477 LOC: variables, personal_token, graphql). Só mudança de caminho de import. Entrega em dias e valida a abordagem.
Construir pkg/registry para resolução de clientes. Substituir os 181 call sites com "api_v4_url" hardcoded.
Consolidar os 40 arquivos de comando quase duplicados (3.901 LOC) em implementações únicas, em ordem de dependência: primeiro os comandos CRUD (variables, personal_token, origin), depois logs/login/whoami/version/completion e, por fim, os dispatchers.
Parametrizar os 19 arquivos da faixa intermediária (2.065 LOC) e consolidar os que reduzirem de forma limpa. Os que resistirem são reclassificados para o Tier 2 e deixados como estão — a reclassificação é um resultado válido, não uma falha. Os cinco dispatchers do grupo (create.go, list.go, delete.go, describe.go, update.go) são a exceção: eles só podem ser consolidados depois que o registry de comandos existir, então entram na fase 2.
Reorganizar o package compartilhado messages/ para que a versão dona de cada string seja explícita, em vez de strings exclusivas de v4 e strings compartilhadas ficarem lado a lado sem marcação.
Critérios de saída: entre 5.000 e 5.900 LOC removidas; toda string em messages/ tem uma versão dona explícita; comportamento idêntico do CLI para clientes de v3 e v4.
Observação: as diferenças de flags no deploy entre v3 e v4 são deliberadas — a v3 está sob política de somente correções de bugs — e permanecem. Não são uma lacuna a ser fechada. A fase 1 apenas torna essa titularidade explícita em messages/, em vez de deixá-la implícita.
Fase 2 — isolamento de recursos do Tier 2 · 8 turnos
Criar pkg/resources/{v3,v4}/ e mover os recursos exclusivos de cada versão (3.161 LOC do lado de v3, mais o conjunto exclusivo de v4).
Substituir setCmds/setV3Cmds pelo registry de dois níveis. A raiz monta sua lista de comandos uma única vez; cada dispatcher preenche seus filhos a partir de registry.Children(verb, f), filtrado por f.APIVersion. É aqui que os cinco dispatchers postergados da fase 1 são consolidados em uma cópia cada, com seus blocos Example gerados a partir dos filhos registrados em vez de hardcoded.
Descontinuar pkg/v3commands/ e pkg/v3api/ como árvores de primeiro nível — o que sobrar delas passa a viver em pkg/resources/v3/.
Critérios de saída: uma única árvore de comandos; diferenças entre versões expressas apenas como entradas de registry; adicionar uma geração passa a significar adicionar um package pkg/resources/vN/.
Fase 3 — extração da pipeline do Tier 3 · 6 turnos
A fase maior e de maior risco; ela mexe no deploy, o caminho crítico voltado ao cliente.
Definir pkg/pipeline — Step, Pipeline, Plan e o runner compartilhado (timing, spinner, logging, dry-run, wrapping de erros, limpeza de órfãos).
Portar a v4 primeiro. Converter manifest.go:80-210 em uma lista declarada de steps sobre os métodos ResourceContext.Apply* existentes. Majoritariamente mecânico — os métodos já têm o formato correto.
Portar a v3 (pkg/v3manifest) para o mesmo runner.
Unificar deploy/deploy_remote/sync sobre o runner. É aqui que deploy_remote/requests.go (462 linhas de diff) e sync/tasks.go (387) são endereçados.
Critérios de saída: os fluxos de deploy das duas versões rodam por um único runner com steps fornecidos por versão; o --dry-run deriva da lista de steps; a v6 requer apenas uma nova lista de steps.
Controle de risco: esta fase merece um rollout gradual (projetos internos → flag opt-in → padrão), em vez de um único release. Orçar de acordo.
Fase 4 — implementação da API v6 · ? (dimensionada quando o schema chegar)
Com as fases 0 a 3 concluídas, a v6 passa a ser:
apiversion.V6 + precedência de resolução (a fase 0 tornou isso um novo case, não uma nova lógica)
Package de cliente pkg/api/v6/
pkg/resources/v6/ para recursos exclusivos de v6
pkg/pipeline/v6/ — a nova lista de steps de deploy
Tipos de manifest/contracts da v6
Nenhuma árvore de comandos nova. Os comandos do Tier 1 funcionam com a v6 através do registry.
13. O que precisamos do time de API (fase 0, trilha A)
Estas são as solicitações formais. Cada uma bloqueia ou remodela a fase 4.
#
Solicitação
Por que o CLI precisa disso
A1
A v6 é aditiva sobre a v4 ou um modelo de recursos substituto?
Determina se os recursos de v4 migram para o Tier 1 (compartilhados com a v6) ou permanecem no Tier 2. Dimensiona diretamente a fase 4.
A2
O modelo de deploy da v6, mesmo em rascunho. O que é criado, em que ordem, com quais dependências?
Isso é a lista de steps do Tier 3. Um rascunho é suficiente; o formato importa mais do que os campos.
A3
Uma mesma conta pode ter direito a mais de uma geração simultaneamente?
Se sim, um único Factory.APIVersion é insuficiente e a resolução passa a ser por comando. Esta é uma premissa estrutural — veja AS-2.
A4
Um sinal de erro distinguível para "esta versão da API não é permitida para esta conta".
O mecanismo de correção pelo servidor no cache de versão depende disso. Sem ele, um cache obsoleto após uma mudança de conta produz falhas confusas em vez de uma nova consulta clara.
A5
Compromisso de estabilidade do schema e um canal de notificação de mudanças.
Determina se a fase 4 pode começar sobre um rascunho ou precisa esperar um schema congelado.
A6
SDK de Go da v6: cronograma, pipeline de geração e se ele substitui o SDK atual.
O sufixo -dev no SDK de v4 atual é, por si só, uma questão de supply chain que vale resolver aqui.
A7
Existe um caminho de migração de v4 para v6 para recursos existentes, e o CLI o expõe?
Um comando migrate seria o único lugar a precisar de duas gerações da API ativas no mesmo processo, algo que o design atual, de uma versão por invocação, não consegue expressar. Barato de acomodar agora, caro depois.
A8
Data de fim de vida da v3, se houver.
A v3 deve sobreviver ao lado da v6 conforme a direção atual. Uma data de descontinuação permitiria apagar pkg/resources/v3/ do Tier 2 por completo e reduzir materialmente o custo de longo prazo.


14. Premissas e registro de decisões
Declaradas explicitamente para que o plano sobreviva à chegada do design real da v6. Cada uma indica o que muda se se provar errada.
#
Premissa
Se estiver errada
AS-1
A v6 é uma substituição completa do modelo de deploy e do vocabulário de recursos (pior caso).
Qualquer coisa menos que isso torna o plano mais barato, não mais caro. Mais recursos de v4 são promovidos ao Tier 1.
AS-2
Uma conta resolve para exatamente uma geração da API por invocação.
Premissa de maior impacto. Se contas puderem ter múltiplas gerações (solicitação A3), Factory.APIVersion passa a ser resolução por comando e o registry ganha uma dimensão. Afeta as fases 0 e 2. Confirmar cedo.
AS-3
A v3 precisa continuar funcional ao lado da v4 e da v6 indefinidamente.
Uma data de descontinuação da v3 (A8) permite apagar pkg/resources/v3/, removendo cerca de 3.161 LOC e uma dimensão do registry.
AS-4
O comportamento dos clientes atuais de v3 e v4 é preservado exatamente ao longo das fases 0 a 3.
Relaxar isso permite que a correção do D3 (falhar de forma explícita em vez de fazer downgrade silencioso para v3) entre já na fase 0, em vez de esperar o próximo major, ao custo da garantia de compatibilidade.
AS-5
A consulta de flags da conta permanece o mecanismo de seleção de versão, generalizada para um enum.
Uma versão declarada pelo projeto (em azion.json) seria uma mudança de UX maior; o trabalho de cache e pré-parsing da fase 0 é pré-requisito de qualquer forma.
AS-6
O trabalho de v6 começa depois da fase 1, em paralelo com as fases 2 e 3.
Se a v6 for antecipada, só a fase 1 já remove a maior parte da redundância e torna até um fork de v6 mais barato. A fase 1 é a jogada que preserva opções.
AS-7
As interfaces do Tier 1 documentam contratos que os packages de API já satisfazem.
Verificado para os 40 arquivos quase duplicados (o diff é a linha de import) e para os 11 arquivos idênticos byte a byte. Risco baixo.


15. Explicitamente fora de escopo
Nomeado para que as fases permaneçam delimitadas.
Consolidação de testes. pkg/v3commands/pkg/v3api/pkg/v3manifest carregam 11.165 LOC de testes majoritariamente duplicados. A consolidação do Tier 1 vai apagar os testes dos arquivos apagados, mas converter o restante em testes table-driven multi-versão está deliberadamente postergado. É o desdobramento natural da fase 1 e deve ser escopado separadamente, quando os limites entre tiers já forem reais. Recomendado como o primeiro trabalho seguinte à fase 1.
Uma suíte de regressão de golden output / end-to-end do CLI. Reduziria materialmente o risco das fases 1 a 3. Excluída aqui para manter a decisão de investimento focada; sinalizada como a adição opcional de maior valor, caso a liderança queira uma rede de segurança mais forte, em especial para a fase 3.
Unificar o vocabulário de recursos de v3 e v4. Um não objetivo deliberado — veja o Tier 2.
Limpeza de UX (por exemplo, a divergência de nomes entre list rule_engine e create rules_engine). Postergada para o próximo major release do CLI.
Reescrever os SDKs gerados. Assunto upstream; acompanhado na solicitação A6.
Design da superfície de comandos da v6. Não pode ser especificado antes de A1 e A2 serem respondidas.


Apêndice — como reproduzir estes números
# Faixas de duplicação entre as duas árvores
for f in $(cd pkg/v3commands && find . -name "*.go" -not -name "*_test.go"); do
  a="pkg/v3commands/$f"; b="pkg/cmd/$f"
  [ -f "$b" ] && printf "%-55s %s\n" "$f" "$(diff "$a" "$b" | grep -c '^[<>]')"
done | sort -k2 -n
 
# Taxa de mudanças espelhadas nos últimos 12 meses
git log --since="12 months ago" --format=%H -- pkg/v3commands pkg/v3api pkg/v3manifest | sort > /tmp/v3c
git log --since="12 months ago" --format=%H -- pkg/cmd pkg/api pkg/manifest       | sort > /tmp/v4c
comm -12 /tmp/v3c /tmp/v4c | wc -l   # as duas árvores
comm -13 /tmp/v3c /tmp/v4c | wc -l   # apenas v4
 
# Packages de v3api idênticos byte a byte
diff -r pkg/v3api/variables pkg/api/variables
diff -r pkg/v3api/personal_token pkg/api/personal_token
 
# O PR espelhado usado como exemplo
git show --stat 2202cb65
Referências principais no código-fonte
Preocupação
Localização
Seleção de versão
pkg/cmd/root/root.go:214-224
Consulta de conta
pkg/cmd/root/utils.go:20
Registro de comandos
pkg/cmd/root/root.go:134-186
Inicialização anterior ao Cobra
cmd/azion/main.go:16-38
Estado global de flags
pkg/cmd/deploy/deploy.go:59-74
Pipeline de apply escrita manualmente
pkg/manifest/manifest.go:80-210
Formato de step de apply já existente
pkg/manifest/resource_applier.go
Precedente de cache com TTL
pkg/cmd/root/pre_command.go:211
Struct de settings
pkg/token/models.go:32-44



