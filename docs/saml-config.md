# Using a SAML Identity Provider

The `management-api` service implements the possibility to connect to an externel identity provider (IDP) service to authenticate for the management interface.

Currently only the SAML protocol is implemented. The `management-api` acts in this case as the Service Provider (SP).

## Configuration

### Preparations

To be able to use the IDP, the following files needs to be available for the service:

- self-signed certificate key pair (e.g., saml.cert and saml.key). To create the self-signed X.509 key pair used by the service provider, you can use something like [(source)](https://github.com/crewjam/saml):

```
openssl req -x509 -newkey rsa:2048 -keyout saml.key -out saml.cert -days 365 -nodes -subj "/CN=myservice.example.com"
```

- HTML template (using go's template engine) that will be displayed upon successful login. One example template is provided in this repository under `templates/saml/login-success.html` .

For example, in a containerised environment (e.g., OpenShift, Kubernetes), mount a folder to the container which includes these files.

### Environment variables

First, we need to turn on this feature:

```
USE_LOGIN_WITH_EXTERNAL_IDP_ENDPOINT=true
```

If this variable is set with true, we need a valid configuration (see following), otherwise the app will crash (panic) at initialisation.
Following environment variables need to be set:

```
SAML_SESSION_CERT_PATH
SAML_SESSION_KEY_PATH
SAML_TEMPLATE_PATH_LOGIN_SUCCESS
SAML_IDP_METADATA_URL
SAML_IDP_URL
SAML_SERVICE_PROVIDER_ROOT_URL
SAML_ENTITY_ID
```

#### How to use the environment variables

##### Certificate:

Configure the certificate and key file path, so that they can be loaded. For example:

```
SAML_SESSION_CERT_PATH="./saml/saml.cert"
SAML_SESSION_KEY_PATH="./saml/saml.key"
```

##### HTML template:

The file path the to HTML template need to point to a valid location on the file system (from where the management-api app is running):

```
SAML_TEMPLATE_PATH_LOGIN_SUCCESS="./saml/login-success.html"
```

##### Add IDP config:

The metadata URL is used to automatically download the XML data file for the IDP configuration:

```
SAML_IDP_METADATA_URL="https://<your-idp-root>/FederationMetadata/2007-06/FederationMetadata.xml"
```

The root URL of the identity provider will be added to the CORS allow origins list. Also this will be used in the log database to identify where the authentication info is coming from:

```
SAML_IDP_URL="https://<your-idp-root>"
```

##### Add SP config:

The root URL should be the address how the service should be called from outside. This is necessary because the server is typically running behind a proxy, and thus cannot know it's own external address.
The root URL of the service provider must contain subroutes (e.g. https://your.server.com/adminapi) as well if relevant.

```
SAML_SERVICE_PROVIDER_ROOT_URL="https://<your-sp-root>"
```

The assertion URL is generated as
`"<your-sp-root>/v1/saml/acs"`.

Entity ID name for the service provider often uses the root URL of the service, but can be configured to contain different values as well, thus it is set separately:


```
SAML_ENTITY_ID="<your-sp-root>"
```
