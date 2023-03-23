package variables

import "k8s.io/apimachinery/pkg/api/resource"

var ClusterRoleName string
var ClusterRoleBindingName string
var StatefulSetReplicas int32 = 1

const LabelKey = "app"
const LabelValue = "openldap"
const RoleBindingServiceAccount = "default"
const SecretName = "openldap-secret"
const ServiceName = "openldap"
const ContainerName = "openldap"
const OpenldapImage = "docker.io/aescanero/openldap-node:0.1"
const StatefulSetName = "openldap"
const ServiceNameHeadLess = "openldap-headless"
const Port int32 = 389
const TargetPort int32 = 1389
const TLSPort int32 = 636
const TargetTLSPort int32 = 1636

var GigaByte = resource.MustParse("1Gi")
var LdifBase = "dc=example"
