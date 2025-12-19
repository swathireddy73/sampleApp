@Library('mySharedLibrary') _

def buildTag = ''

pipeline {
    agent { label 'build-agent' }

    parameters {
        string(name: 'APP_VERSION', defaultValue: 'v1', description: 'App version to build and deploy')
        choice(name: 'ENV', choices: ['dev','staging','prod'], description: 'Target environment')
    }

    environment {
        HELM_RELEASE = 'userapp-release'
        K8S_NAMESPACE = "${params.ENV}"
        SONAR_PROJECT_KEY = 'sampleapp'
        SONAR_HOST_URL = 'http://20.75.196.235:9000/'
    }

    stages {

        stage('Generate Tag') {
            steps {
                script {
                    buildTag = generateTag()
                    echo "Generated build tag: ${buildTag}"
                }
            }
        }

        stage('Checkout Code') {
            steps {
                script {
                    def branchToBuild = params.BRANCH ?: 'master'
                    git branch: branchToBuild,
                        url: 'https://github.com/swathireddy73/sampleApp.git',
                        credentialsId: '40a1d4f8-1be4-4f42-a7f1-a4da2eb75b93'
                }
            }
        }

        stage('SonarQube Analysis') {
            steps {
                script {
                    def scannerHome = tool name: 'mysonarscanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
                    withSonarQubeEnv('sonarkube-swathi') {
                        sh """
                            ${scannerHome}/bin/sonar-scanner \
                              -Dsonar.projectKey=${env.SONAR_PROJECT_KEY} \
                              -Dsonar.sources=. \
                              -Dsonar.host.url=${env.SONAR_HOST_URL} \
                              -Dsonar.login=$SONAR_AUTH_TOKEN
                        """
                    }
                }
            }
        }

        stage('Quality Gate') {
            steps {
                timeout(time: 5, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: false
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                sh "docker build -t swathireddy73/sampleapp:${params.APP_VERSION} ."
            }
        }

        stage('Push Docker Image') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'b13e918c-c5ee-412e-9dc2-75bf2eabeec3',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh """
                        echo \$DOCKER_PASS | docker login -u \$DOCKER_USER --password-stdin
                        docker push swathireddy73/sampleapp:${params.APP_VERSION}
                    """
                }
            }
        }

      stage('Azure Login & AKS Setup') {
    steps {
        withCredentials([usernamePassword(
            credentialsId: 'aks-login',
            usernameVariable: 'AZURE_CLIENT_ID',
            passwordVariable: 'AZURE_CLIENT_SECRET'
        )]) {
            sh '''
                az login --service-principal -u "$AZURE_CLIENT_ID" -p "$AZURE_CLIENT_SECRET" --tenant 2b32b1fa-7899-482e-a6de-be99c0ff5516
                az aks get-credentials --resource-group rg-dev-flux --name aks-dev-flux-cluster --overwrite-existing
                kubectl get pods -n default
            '''
        }
    }
}



        stage('Deploy with Helm') {
            steps {
                sh """
                    echo "Deploying Helm chart to AKS..."
                    set -e
                    helm upgrade --install ${env.HELM_RELEASE} ./helm-chart \
                        --namespace ${params.ENV} \
                        --set image.tag=${params.APP_VERSION} \
                        --create-namespace
                    kubectl rollout status deployment/${env.HELM_RELEASE} -n ${params.ENV}
                """
            }
        }

    } // stages
}
