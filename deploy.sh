#!/bin/bash
cp index.html docker-deploy/web/
cp admin_index.html docker-deploy/web/
cp -r dist docker-deploy/web/
cp -r core docker-deploy/
cp -r static docker-deploy/web/
cp -r config docker-deploy/web/
cp ahpuoj docker-deploy/web/ahpuoj
