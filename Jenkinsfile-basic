pipeline {
    agent { label 'build-agent' }

    stages {
        stage('Checkout Code') {
            steps {
                git url: 'https://github.com/gititc778/sampleApp.git', branch: 'master'
            }
        }

        stage('Build Docker Image') {
            steps {
                sh 'docker build -t sampleapp:${BUILD_NUMBER} .'
            }
        }

        stage('Push to Docker Registry') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'docker-credentials', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    sh '''
                        echo "${DOCKER_PASS}" | docker login -u "${DOCKER_USER}" --password-stdin
                        docker tag sampleapp:${BUILD_NUMBER} ${DOCKER_USER}/sampleapp:${BUILD_NUMBER}
                        docker push ${DOCKER_USER}/sampleapp:${BUILD_NUMBER}
                    '''
                }
            }
        }

        stage('Deploy to Kubernetes'){
            steps {
                withCredentials([file(credentialsId: 'kubeconfig-creds', variable: 'KUBECONFIG')]) {
                sh '''
                    sed -i "s/IMAGE_TAG/${BUILD_NUMBER}/g" deployment.yaml
                    kubectl apply -f deployment.yaml
                '''
                }
            }
        }
    }
}
