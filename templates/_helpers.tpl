{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "awx-helm.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "awx-helm.fullname" -}}
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
Create chart name and version as used by the chart label.
*/}}
{{- define "awx-helm.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{/* Create the template for the environment.sh script */}}
{{- define "awx-helm.environment_sh" -}}
DEFAULT_LOCAL_TMP="/var/lib/awx/.ansible/tmp"
DATABASE_USER={{ .Values.postgresql.postgresqlUsername }}
DATABASE_NAME={{ .Values.postgresql.postgresqlDatabase }}
DATABASE_HOST={{ .Values.postgresql.postgresqlHostname }}
DATABASE_PORT={{ .Values.postgresql.postgresqlPort }}
DATABASE_PASSWORD={{ .Values.postgresql.postgresqlPassword }}
MEMCACHED_HOST={{ default "localhost" .Values.memcached.memcachedHostname }}
MEMCACHED_PORT={{ .Values.memcached.memcachedPort }}
RABBITMQ_HOST={{ default "localhost" .Values.rabbitmq.rabbitmqHostname }}
RABBITMQ_PORT="{{ default 5672 .Values.rabbitmq.rabbitmqPort }}"
AWX_ADMIN_USER={{ .Values.adminUser }}
AWX_ADMIN_PASSWORD={{ .Values.adminPassword }}
{{- end -}}

{{/* Create the template for the credentials.py script */}}
{{- define "awx-helm.credentials_py" -}}
DATABASES = {
    'default': {
        'ATOMIC_REQUESTS': True,
        'ENGINE': 'awx.main.db.profiled_pg',
        'NAME': "{{ .Values.postgresql.postgresqlDatabase }}",
        'USER': "{{ .Values.postgresql.postgresqlUsername }}",
        'PASSWORD': "{{ .Values.postgresql.postgresqlPassword }}",
        'HOST': "{{ .Values.postgresql.postgresqlHostname }}",
        'PORT': "{{ .Values.postgresql.postgresqlPort }}",
    }
}

BROKER_URL = 'amqp://{}:{}@{}:{}/{}'.format(
        "{{ .Values.rabbitmq.rabbitmqUsername }}",
        "{{ .Values.rabbitmq.rabbitmqPassword }}",
        "localhost",
        "{{ default "5672" .Values.rabbitmq.rabbitmqPort }}",
        "{{ default "awx" .Values.rabbitmq.rabbitmqVhost }}")

CHANNEL_LAYERS = {
        'default': {'BACKEND': 'asgi_amqp.AMQPChannelLayer',
                    'ROUTING': 'awx.main.routing.channel_routing',
                    'CONFIG': {'url': BROKER_URL}}
}

CHANNEL_LAYERS = {
    'default': {'BACKEND': 'asgi_amqp.AMQPChannelLayer',
                'ROUTING': 'awx.main.routing.channel_routing',
                'CONFIG': {'url': BROKER_URL}}
}
{{- end -}}
