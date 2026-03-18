{{/*  
Expand the name of the chart.  
*/}}  
{{- define "pizza-hub.name" -}}  
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}  
{{- end }}  

{{/*  
Create a default fully qualified app name.  
*/}}  
{{- define "pizza-hub.fullname" -}}  
{{- if .Values.fullnameOverride }}  
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}  
{{- else }}  
{{- $name := default .Chart.Name .Values.nameOverride }}  
{{- if contains $name .Release.Name }}  
{{- .Release.Name | trunc 63 | trimSuffix "-" }}  
{{- else }}  
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}  
{{- end }}  
{{- end }}  
{{- end }}  

{{/*  
Create chart name and version as used by the chart label.  
*/}}  
{{- define "pizza-hub.chart" -}}  
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}  
{{- end }}  

{{/*  
Common labels  
*/}}  
{{- define "pizza-hub.labels" -}}  
helm.sh/chart: {{ include "pizza-hub.chart" . }}  
{{ include "pizza-hub.selectorLabels" . }}  
{{- if .Chart.AppVersion }}  
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}  
{{- end }}  
app.kubernetes.io/managed-by: {{ .Release.Service }}  
{{- end }}  

{{/*  
Selector labels  
*/}}  
{{- define "pizza-hub.selectorLabels" -}}  
app.kubernetes.io/name: {{ include "pizza-hub.name" . }}  
app.kubernetes.io/instance: {{ .Release.Name }}  
{{- end }}  

{{/*  
Create the name of the service account to use  
*/}}  
{{- define "pizza-hub.serviceAccountName" -}}  
{{- if .Values.serviceAccount.create }}  
{{- default (include "pizza-hub.fullname" .) .Values.serviceAccount.name }}  
{{- else }}  
{{- default "default" .Values.serviceAccount.name }}  
{{- end }}  
{{- end }}  