#!/bin/bash

# Stop and remove all Docker containers
echo "Stopping all Docker containers..."
docker stop $(docker ps -a -q)

echo "Removing all Docker containers..."
docker rm $(docker ps -a -q)

# Prune and remove Docker volumes
echo "Pruning Docker volumes..."
docker volume prune -f

echo "Listing Docker volumes..."
docker volume ls

echo "Removing all Docker volumes..."
docker volume rm $(docker volume ls -q)

echo "Listing Docker volumes again to confirm removal..."
docker volume ls

# Prune Docker networks
echo "Pruning Docker networks..."
docker network prune -f

# Remove old artifacts and create necessary directories
echo "Removing old organizations and channel artifacts..."
sudo rm -rf organizations/peerOrganizations
sudo rm -rf organizations/ordererOrganizations
sudo rm -rf channel-artifacts/
mkdir channel-artifacts

# Set environment variables for Fabric binaries and configuration path
echo "Setting environment variables for Fabric binaries and configuration..."
export PATH=${PWD}/../fabric-samples/bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/../fabric-samples/config

# Start the Certificate Authority
echo "Starting the Certificate Authority..."
docker-compose -f docker/docker-compose-ca.yaml up -d

# Wait for CA to start
echo "Waiting for CA to start..."
sleep 10

# Execute registerEnroll.sh script to create organizations
echo "Creating Org1...Madrid"
. ./organizations/fabric-ca/registerEnroll.sh && createMadrid
sleep 5

echo "Creating Org2...Bogota"
. ./organizations/fabric-ca/registerEnroll.sh && createBogota
sleep 5

echo "Creating Orderer..."
. ./organizations/fabric-ca/registerEnroll.sh && createOrderer
sleep 5

# Start the network (peers, orderer, and CouchDB)
echo "Starting the network (peers, orderer, and CouchDB)..."
docker-compose -f docker/docker-compose-test-net.yaml up -d

# Wait for network to start
echo "Waiting for network to start..."
sleep 20

# Display the status of the Docker containers
echo "Displaying Docker containers status..."
docker ps -a

# Display logs for peer0.madrid.universidades.com
echo "Displaying logs for peer0.madrid.universidades.com..."
docker logs peer0.madrid.universidades.com

# Generate channel creation transaction
echo "Generating channel creation transaction..."
export FABRIC_CFG_PATH=${PWD}/configtx
configtxgen -profile TwoOrgsApplicationGenesis -outputBlock ./channel-artifacts/universidadeschannel.block -channelID universidadeschannel

# Set environment variables for Orderer
echo "Setting environment variables for Orderer..."
export FABRIC_CFG_PATH=${PWD}/../fabric-samples/config
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/universidades.com/orderers/orderer.universidades.com/msp/tlscacerts/tlsca.universidades.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/universidades.com/orderers/orderer.universidades.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/universidades.com/orderers/orderer.universidades.com/tls/server.key

# Wait for channel creation transaction to be generated
echo "Waiting for channel creation transaction to be generated..."
sleep 5

# Join the channel
echo "Orderer joining the channel..."
osnadmin channel join --channelID universidadeschannel --config-block ./channel-artifacts/universidadeschannel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
sleep 5

# List channels
echo "Listing channels for the Orderer..."
osnadmin channel list -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
sleep 5

# Set environment variables for peer0.madrid.universidades.com

export CORE_PEER_TLS_ENABLED=true
export PEER0_MADRID_CA=${PWD}/organizations/peerOrganizations/madrid.universidades.com/peers/peer0.madrid.universidades.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="MadridMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_MADRID_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/madrid.universidades.com/users/Admin@madrid.universidades.com/msp
export CORE_PEER_ADDRESS=localhost:7051

# Wait for environment variables to be set
echo "Waiting for environment variables to be set..."
sleep 5

# Peer0.madrid.universidades.com joins the channel
echo "peer0.madrid.universidades.com joining the channel..."
peer channel join -b ./channel-artifacts/universidadeschannel.block
sleep 5

# Set environment variables for peer0.bogota.universidades.com
echo "Setting environment variables for peer0.bogota.universidades.com..."
export PEER0_BOGOTA_CA=${PWD}/organizations/peerOrganizations/bogota.universidades.com/peers/peer0.bogota.universidades.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="BogotaMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_BOGOTA_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/bogota.universidades.com/users/Admin@bogota.universidades.com/msp
export CORE_PEER_ADDRESS=localhost:9051


# Wait for environment variables to be set
echo "Waiting for environment variables to be set..."
sleep 5

# Peer0.bogota.universidades.com joins the channel
echo "peer0.bogota.universidades.com joining the channel..."
peer channel join -b ./channel-artifacts/universidadeschannel.block
sleep 5

# Install and approve chaincode on peer0.madrid.universidades.com
echo "Packaging chaincode..."
peer lifecycle chaincode package testt.tar.gz --path chaincodes/test/go --lang golang --label testt_1.0
sleep 2

echo "Setting environment variables for peer0.madrid.universidades.com install..."
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="MadridMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/madrid.universidades.com/peers/peer0.madrid.universidades.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/madrid.universidades.com/users/Admin@madrid.universidades.com/msp
export CORE_PEER_ADDRESS=localhost:7051
echo "Installing chaincode on peer0.madrid.universidades.com..."
peer lifecycle chaincode install testt.tar.gz
sleep 2

echo "Setting environment variables for peer0.bogota.universidades.com install..."
export CORE_PEER_LOCALMSPID="BogotaMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/bogota.universidades.com/peers/peer0.bogota.universidades.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/bogota.universidades.com/users/Admin@bogota.universidades.com/msp
export CORE_PEER_ADDRESS=localhost:9051
echo "Installing chaincode on peer0.bogota.universidades.com..."
peer lifecycle chaincode install testt.tar.gz
sleep 2

echo "Querying installed chaincode to get package ID..."
CC_PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep testt_1.0 | sed -n 's/^Package ID: //;s/, Label:.*$//;p')
echo "Chaincode package ID is ${CC_PACKAGE_ID}"
sleep 2

echo "Approving chaincode on peer0.bogota.universidades.com..."
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.universidades.com --channelID universidadeschannel --signature-policy "OR('MadridMSP.member','BogotaMSP.member')" --name testt --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/universidades.com/orderers/orderer.universidades.com/msp/tlscacerts/tlsca.universidades.com-cert.pem
sleep 2

echo "Setting environment variables for peer0.madrid.universidades.com approval..."
export CORE_PEER_LOCALMSPID="MadridMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/madrid.universidades.com/peers/peer0.madrid.universidades.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/madrid.universidades.com/users/Admin@madrid.universidades.com/msp
export CORE_PEER_ADDRESS=localhost:7051

echo "Approving chaincode on peer0.madrid.universidades.com..."
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.universidades.com --channelID universidadeschannel --signature-policy "OR('MadridMSP.member','BogotaMSP.member')" --name testt --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/universidades.com/orderers/orderer.universidades.com/msp/tlscacerts/tlsca.universidades.com-cert.pem
sleep 2

echo "Checking commit readiness for chaincode..."
peer lifecycle chaincode checkcommitreadiness --channelID universidadeschannel --name testt --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/universidades.com/orderers/orderer.universidades.com/msp/tlscacerts/tlsca.universidades.com-cert.pem --output json
sleep 2

echo "Committing chaincode definition on the channel..."
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.universidades.com --signature-policy "OR('MadridMSP.member','BogotaMSP.member')" --channelID universidadeschannel --name testt --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/universidades.com/orderers/orderer.universidades.com/msp/tlscacerts/tlsca.universidades.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/madrid.universidades.com/peers/peer0.madrid.universidades.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/bogota.universidades.com/peers/peer0.bogota.universidades.com/tls/ca.crt
sleep 2

echo "Querying committed chaincode..."
peer lifecycle chaincode querycommitted --channelID universidadeschannel --name testt --cafile ${PWD}/organizations/ordererOrganizations/universidades.com/orderers/orderer.universidades.com/msp/tlscacerts/tlsca.universidades.com-cert.pem
sleep 2

echo "Invoking chaincode to initialize ledger..."
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.universidades.com --tls --cafile ${PWD}/organizations/ordererOrganizations/universidades.com/orderers/orderer.universidades.com/msp/tlscacerts/tlsca.universidades.com-cert.pem -C universidadeschannel -n testt --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/madrid.universidades.com/peers/peer0.madrid.universidades.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'
sleep 5

echo "Querying chaincode to ObtenerTodosLosAlumnos..."
peer chaincode query -C universidadeschannel -n testt -c '{"Args":["ObtenerTodosLosAlumnos"]}'
sleep 2

# Crear un nuevo alumno
echo "Invoking chaincode to CrearAlumno...[3]"
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.universidades.com --tls --cafile ${PWD}/organizations/ordererOrganizations/universidades.com/orderers/orderer.universidades.com/msp/tlscacerts/tlsca.universidades.com-cert.pem -C universidadeschannel -n testt --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/madrid.universidades.com/peers/peer0.madrid.universidades.com/tls/ca.crt -c '{"function":"CrearAlumno","Args":["3","Pedro","Martinez","30","Psicologia"]}'
sleep 5

# Validar que el alumno ha sido registrado consultando sus detalles
echo "Querying chaincode to LeerAlumno...[3]"
peer chaincode query -C universidadeschannel -n testt -c '{"Args":["LeerAlumno","3"]}'
sleep 2


# Eliminar el alumno con ID 2
echo "Invoking chaincode to EliminarAlumno...[2]"
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.universidades.com --tls --cafile ${PWD}/organizations/ordererOrganizations/universidades.com/orderers/orderer.universidades.com/msp/tlscacerts/tlsca.universidades.com-cert.pem -C universidadeschannel -n testt --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/madrid.universidades.com/peers/peer0.madrid.universidades.com/tls/ca.crt -c '{"function":"EliminarAlumno","Args":["2"]}'
sleep 5

# Validar que el alumno ha sido eliminado
echo "Querying chaincode to validate EliminarAlumno..."
peer chaincode query -C universidadeschannel -n testt -c '{"Args":["LeerAlumno","2"]}'
sleep 2


echo "Querying chaincode to ObtenerTodosLosAlumnos..."
peer chaincode query -C universidadeschannel -n testt -c '{"Args":["ObtenerTodosLosAlumnos"]}'
sleep 2

echo "Chaincode operations completed successfully."

