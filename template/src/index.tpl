{% extends "internal/layout.tpl" %}
{% block content %}
<h1>Index</h1>
{% if context.a %}
<h2>
	{{ context.a }}
</h2>
{% endif %}
{% endblock %}