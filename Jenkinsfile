@Library('mySharedLibrary') _

def buildTag = ''

pipeline {
    agent { label 'build-agent' }

    parameters {
        string(name: 'BRANCH', defaultValue: 'master', description: 'Branch to build')
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
            // Add Git credentials here
            def branchToBuild = params.BRANCH ?: 'master'

            // Checkout the repo with credentials
            git branch: branchToBuild,
                url: 'https://github.com/swathireddy73/sampleApp.git',
                credentialsId: '40a1d4f8-1be4-4f42-a7f1-a4da2eb75b93'  // replace with your Jenkins credential ID
        }
    }
}
stage('SonarQube Analysis') {
    steps {
        withSonarQubeEnv('sonarkube-swathipothula') {
            sh '''
                sonar-scanner \
                  -Dsonar.projectKey=sampleapp \
                  -Dsonar.sources=. \
                  -Dsonar.host.url=http://20.75.196.235:9000 \
                  -Dsonar.login=$SONAR_AUTH_TOKEN
            '''
        }
    }
}


        stage('Quality Gate') {
            steps {
                timeout(time: 5, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                sh "docker build -t sampleapp:${params.APP_VERSION} ."
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
                        docker login -u $DOCKER_USER -p $DOCKER_PASS
                        docker tag sampleapp:${params.APP_VERSION} mydockerhubuser/sampleapp:${params.APP_VERSION}
                        docker push mydockerhubuser/sampleapp:${params.APP_VERSION}
                    """
                }
            }
        }

        stage('Deploy with Helm') {
            steps {
                withCredentials([file(credentialsId: 'kubeconfig-file', variable: 'KUBECONFIG')]) {
                    sh """
                        helm upgrade --install ${HELM_RELEASE} ./helm-chart \
                        --namespace ${K8S_NAMESPACE} \
                        --set image.tag=${params.APP_VERSION}
                    """
                }
            }
        }
    }
}
