{{- define "calendar-chart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "calendar-chart.fullname" -}}
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

{{- define "calendar-chart.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "calendar-chart.labels" -}}
helm.sh/chart: {{ include "calendar-chart.chart" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "calendar-chart.selectorLabels" -}}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "calendar-chart.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "calendar-chart.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "calendar-chart.postgresSecret" -}}
{{- if .Values.postgresql.existingSecret }}
{{- .Values.postgresql.existingSecret }}
{{- else }}
{{- include "calendar-chart.fullname" . }}-postgres
{{- end }}
{{- end }}

{{- define "calendar-chart.rabbitmqSecret" -}}
{{- if .Values.rabbitmq.existingSecret }}
{{- .Values.rabbitmq.existingSecret }}
{{- else }}
{{- include "calendar-chart.fullname" . }}-rabbitmq
{{- end }}
{{- end }}

{{- define "calendar-chart.postgresUrl" -}}
{{- printf "postgres://%s@%s:%d/%s?sslmode=disable" .Values.postgresql.user .Values.postgresql.host (int .Values.postgresql.port) .Values.postgresql.database }}
{{- end }}

{{- define "calendar-chart.rabbitmqUrl" -}}
{{- printf "amqp://%s@%s:%d/" .Values.rabbitmq.user .Values.rabbitmq.host (int .Values.rabbitmq.port) }}
{{- end }}
