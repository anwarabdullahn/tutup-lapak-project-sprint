{{- define "tutup-lapak.name" -}}
{{ .Chart.Name }}
{{- end -}}

{{- define "tutup-lapak.fullname" -}}
{{- printf "%s" .Chart.Name -}}
{{- end -}}

{{- define "tutup-lapak.image" -}}
{{- $registry := .Values.global.image.registry -}}
{{- $tag := .Values.global.image.tag -}}
{{ printf "%s/%s:%s" $registry .imageName $tag }}
{{- end -}}
