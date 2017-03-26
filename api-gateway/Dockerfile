FROM scratch
ADD server/api-gateway-docker /
COPY client/public .
EXPOSE 8080
CMD ["/api-gateway-docker", "-d", "."]
