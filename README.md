## Twitter Operator for Kubernetes

Este es un operador para manejar tweets usando Kubernetes.

Es una idea bastante estúpida, pero funciona, cuando el operador inicia, consume la lista de tweets del usuario, y carga esos tweets como "Tweets" en Kubernetes (por defecto sólo carga los últimos 20 tweets)

Luego podemos crear un nuevo tweet aplicando el ejemplo que está en `config/samples`

### GOTCHAs

Cuando el controlador inicia, borra todos los tweets que están cargados en Kubernetes (no de twitter) ... esto lo hacemos ya que tenemos que refrescar la lista de tweets todas las veces que iniciamos el controlador, ya que el controlador solo sincroniza los tweets cuando hay algun cambio del lado de Kubernetes (no de Twitter) ... para poder lograr lo segundo, deberiamos crear un server que recibe los tweets de Twitter (por medio de un webhook) para poder cargar los nuevos tweets cada vez que se cree uno desde Twitter.

Este operador fue usado para una demo en [Nerdear.la 2019](https://nerdear.la)
