pipeline {
    agent { docker { image 'golang' } }
    stages {
        stage('build') {
            steps {
                go version
            }
        }
    }
}