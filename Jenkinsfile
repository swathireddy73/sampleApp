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
                              -Dsonar.projectKey=sampleapp \
                              -Dsonar.sources=. \
                              -Dsonar.host.url=http://20.75.196.235:9000 \
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

        stage('Deploy with Helm') {
            steps {
                withCredentials([file(credentialsId: 'kubeconfig-file', variable: 'KUBECONFIG')]) {
                    sh """
                        echo "Testing Kubernetes connection..."
                        kubectl get nodes
                        
                        echo "Deploying Helm chart..."
                        helm upgrade --install ${HELM_RELEASE} ./helm-chart \
                          --namespace ${params.ENV} \
                          --set image.tag=${params.APP_VERSION}
                    """
                }
            }
        }
    }
}
