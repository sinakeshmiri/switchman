## Switchman
Switchman uses different DNS providers' APIs and updates A records in order to redirect clients whenever a host goes down.

This is useful in case you have multiple LoadBalancers and you want to distribute traffic between your LoadBalancers. 

Switchman continuously healthchecks the upstreams and update A records accordingly.

## Supported dns APIs
- Cloudflare 

## Usage
Fill variables in config/config.json

## TODO
- [ ] k8s integration
- [ ] better cmdline options (env vars, logging options, ...)
- [ ] add tests
- [ ] add more dns APIs
- [ ] README: how to deploy LoadBalancers in HA mode
- [ ] README: contribute section
- [ ] README: usage section
