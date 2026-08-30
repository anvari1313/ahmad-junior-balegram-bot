{{/*
Common labels applied to all resources.
*/}}
{{- define "bot.labels" -}}
helm.sh/chart: {{ include "bot.chart" . }}
{{ include "bot.selectorLabels" . }}
{{- with .Chart.AppVersion }}
app.kubernetes.io/version: {{ . | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels (immutable for the lifetime of a StatefulSet).
*/}}
{{- define "bot.selectorLabels" -}}
app.kubernetes.io/name: {{ include "bot.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Chart name and version as a single label.
*/}}
{{- define "bot.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Base resource name: fullname override, else chart name (stable across releases).
Used for selector app.kubernetes.io/name so it does not change with release.
*/}}
{{- define "bot.name" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/*
Fully-qualified resource name: release + chart name (or override).
*/}}
{{- define "bot.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
All application env vars (one entry per set variable); empty when nothing
is configured, so the env block can be omitted entirely.
*/}}
{{- define "bot.envVars" -}}
{{- $entries := list
      (include "bot.envVar" (list "TELEGRAM_BOT_TOKEN" .Values.config.telegramBotToken))
      (include "bot.envVar" (list "TELEGRAM_API_URL" .Values.config.telegramApiUrl))
      (include "bot.envVar" (list "ZAI_API_KEY" .Values.config.zaiApiKey))
      (include "bot.envVar" (list "ZAI_BASE_URL" .Values.config.zaiBaseUrl))
      (include "bot.envVar" (list "ZAI_MODEL" .Values.config.zaiModel)) -}}
{{- $out := "" -}}
{{- range $e := $entries -}}
{{- if $e -}}
{{- $out = printf "%s%s\n" $out (trim $e) -}}
{{- end -}}
{{- end -}}
{{- trim $out -}}
{{- end -}}

{{/*
ServiceAccount name.
*/}}
{{- define "bot.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "bot.fullname" .) .Values.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{/*
Image reference: repository:tag, falling back to chart appVersion.
*/}}
{{- define "bot.image" -}}
{{- $tag := .Values.image.tag | default .Chart.AppVersion -}}
{{- printf "%s:%s" .Values.image.repository $tag -}}
{{- end -}}

{{/*
Renders one env var entry: a literal `value` when set, else a `secretKeyRef`
when `secretRef.name`+`secretRef.key` are set, else nothing (app default).
Usage: (include "bot.envVar" (list "VAR_NAME" .Values.config.someVar))
*/}}
{{- define "bot.envVar" -}}
{{- $name := index . 0 -}}
{{- $spec := index . 1 -}}
{{- if $spec.value -}}
- name: {{ $name | quote }}
  value: {{ $spec.value | quote }}
{{- else if and $spec.secretRef $spec.secretRef.name $spec.secretRef.key -}}
- name: {{ $name | quote }}
  valueFrom:
    secretKeyRef:
      name: {{ $spec.secretRef.name | quote }}
      key: {{ $spec.secretRef.key | quote }}
{{- end }}
{{- end -}}
