#!/bin/bash

function createEmpresa() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/empresa.certificacion.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/empresa.certificacion.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname ca-empresa --tls.certfiles "${PWD}/organizations/fabric-ca/empresa/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-empresa.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-empresa.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-empresa.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-empresa.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/msp/config.yaml"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-empresa --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/empresa/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering peer1"
  set -x
  fabric-ca-client register --caname ca-empresa --id.name peer1 --id.secret peer1pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/empresa/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-empresa --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/empresa/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-empresa --id.name empresaadmin --id.secret empresaadminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/empresa/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-empresa -M "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer0.empresa.certificacion.com/msp" --csr.hosts peer0.empresa.certificacion.com --tls.certfiles "${PWD}/organizations/fabric-ca/empresa/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer0.empresa.certificacion.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-empresa -M "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer0.empresa.certificacion.com/tls" --enrollment.profile tls --csr.hosts peer0.empresa.certificacion.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/empresa/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer0.empresa.certificacion.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer0.empresa.certificacion.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer0.empresa.certificacion.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer0.empresa.certificacion.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer0.empresa.certificacion.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer0.empresa.certificacion.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer0.empresa.certificacion.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer0.empresa.certificacion.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/tlsca/tlsca.empresa.certificacion.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/ca"
  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer0.empresa.certificacion.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/ca/ca.empresa.certificacion.com-cert.pem"

  infoln "Generating the peer1 msp"
  set -x
  fabric-ca-client enroll -u https://peer1:peer1pw@localhost:7054 --caname ca-empresa -M "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer1.empresa.certificacion.com/msp" --csr.hosts peer1.empresa.certificacion.com --tls.certfiles "${PWD}/organizations/fabric-ca/empresa/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer1.empresa.certificacion.com/msp/config.yaml"

  infoln "Generating the peer1-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer1:peer1pw@localhost:7054 --caname ca-empresa -M "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer1.empresa.certificacion.com/tls" --enrollment.profile tls --csr.hosts peer1.empresa.certificacion.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/empresa/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer1.empresa.certificacion.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer1.empresa.certificacion.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer1.empresa.certificacion.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer1.empresa.certificacion.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer1.empresa.certificacion.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer1.empresa.certificacion.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer1.empresa.certificacion.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer1.empresa.certificacion.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/tlsca/tlsca.empresa.certificacion.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/ca"
  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/peers/peer1.empresa.certificacion.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/ca/ca.empresa.certificacion.com-cert.pem"


  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca-empresa -M "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/users/User1@empresa.certificacion.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/empresa/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/users/User1@empresa.certificacion.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://empresaadmin:empresaadminpw@localhost:7054 --caname ca-empresa -M "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/users/Admin@empresa.certificacion.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/empresa/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/empresa.certificacion.com/users/Admin@empresa.certificacion.com/msp/config.yaml"
}

function createOrderer() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/ordererOrganizations/certificacion.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/certificacion.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca-orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/ordererOrganizations/certificacion.com/msp/config.yaml"

  infoln "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the orderer admin"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/msp" --csr.hosts orderer.certificacion.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/certificacion.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/tls" --enrollment.profile tls --csr.hosts orderer.certificacion.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/tls/ca.crt"
  cp "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/tls/signcerts/"* "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/tls/server.crt"
  cp "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/tls/keystore/"* "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/tls/server.key"

  mkdir -p "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/msp/tlscacerts/tlsca.certificacion.com-cert.pem"

  mkdir -p "${PWD}/organizations/ordererOrganizations/certificacion.com/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/certificacion.com/orderers/orderer.certificacion.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/certificacion.com/msp/tlscacerts/tlsca.certificacion.com-cert.pem"

  infoln "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/certificacion.com/users/Admin@certificacion.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/certificacion.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/certificacion.com/users/Admin@certificacion.com/msp/config.yaml"
}

