def buildTag = ''

pipeline {
    agent { label 'build-agent' }

    stages {
        stage('Generate Tag') {
            steps {
                script {
                    def date = new Date().format('yyyyMMdd')
                    buildTag = "${date}.${env.BUILD_NUMBER}"
                    currentBuild.displayName = buildTag

                   
                    sh "echo BUILD_TAG=${buildTag} > build.env"
                }
            }
        }

        stage('Use Tag') {
            steps {
                script {
                    echo "The build tag is: ${buildTag}"
                }
            }
        }

        stage('Checkout Code') {
            steps {
                git url: 'https://github.com/gititc778/sampleApp.git', branch: 'master'
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    sh "docker build -t sampleapp:${buildTag} ."
                }
            }
        }

        stage('Push to Docker Registry') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'docker-credentials', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    script {
                        sh """
                            echo "${DOCKER_PASS}" | docker login -u "${DOCKER_USER}" --password-stdin
                            docker tag sampleapp:${buildTag} ${DOCKER_USER}/sampleapp:${buildTag}
                            docker push ${DOCKER_USER}/sampleapp:${buildTag}
                        """
                    }
                }
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                script {
                    input message: "Do you want to proceed with Kubernetes deployment?", ok: 'Deploy'
                }

                withCredentials([file(credentialsId: 'kubeconfig-creds', variable: 'KUBECONFIG')]) {
                    script {
                        sh """
                            sed -i "s/IMAGE_TAG/${buildTag}/g" deployment.yaml
                            kubectl apply -f deployment.yaml
                        """
                    }
                }
            }
        }
    }
}
