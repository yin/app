properties([buildDiscarder(logRotator(numToKeepStr: '20'))])

pipeline {
    agent {
        label 'linux && x86_64'
    }

    options {
        skipDefaultCheckout(true)
    }

    stages {
        stage('Build') {
            parallel {
                stage("Validate") {
                    agent {
                        label 'ubuntu-1604-aufs-edge'
                    }
                    steps {
                        dir('src/github.com/docker/app') {
                            checkout scm
                            sh 'make -f docker.Makefile lint'
                            sh 'make -f docker.Makefile check-vendor'
                        }
                    }
                    post {
                        always {
                            deleteDir()
                        }
                    }
                }
                stage("Binaries"){
                    agent {
                        label 'ubuntu-1604-aufs-edge'
                    }
                    steps  {
                        dir('src/github.com/docker/app') {
                            script {
                                try {
                                    checkout scm
                                    sh 'make -f docker.Makefile cli-cross cross e2e-cross tars'
                                    dir('bin') {
                                        stash name: 'binaries'
                                    }
                                    dir('e2e') {
                                        stash name: 'e2e'
                                    }
                                    dir('examples') {
                                        stash name: 'examples'
                                    }
                                    if(!(env.BRANCH_NAME ==~ "PR-\\d+")) {
                                        stash name: 'artifacts', includes: 'bin/*.tar.gz', excludes: 'bin/*-e2e-*'
                                        archiveArtifacts 'bin/*.tar.gz'
                                    }
                                } finally {
                                    def clean_images = /docker image ls --format="{{.Repository}}:{{.Tag}}" '*$BUILD_TAG*' | xargs docker image rm -f/
                                    sh clean_images
                                }
                            }
                        }
                    }
                    post {
                        always {
                            deleteDir()
                        }
                    }
                }
                stage('Build Invocation image'){
                    agent {
                        label 'ubuntu-1604-aufs-edge'
                    }
                    steps {
                        dir('src/github.com/docker/app') {
                            checkout scm
                            sh 'make -f docker.Makefile save-invocation-image'
                            sh 'make -f docker.Makefile INVOCATION_IMAGE_TAG=$BUILD_TAG-coverage OUTPUT=coverage-invocation-image.tar save-invocation-image-tag'
                            sh 'make -f docker.Makefile INVOCATION_IMAGE_TAG=$BUILD_TAG-coverage-experimental OUTPUT=coverage-experimental-invocation-image.tar save-invocation-image-tag'
                            dir('_build') {
                                stash name: 'invocation-image', includes: 'invocation-image.tar'
                                stash name: 'coverage-invocation-image', includes: 'coverage-invocation-image.tar'
                                stash name: 'coverage-experimental-invocation-image', includes: 'coverage-experimental-invocation-image.tar'
                            }
                        }
                    }
                    post {
                        always {
                            dir('src/github.com/docker/app') {
                                sh 'docker rmi docker/cnab-app-base:$BUILD_TAG'
                                sh 'docker rmi docker/cnab-app-base:$BUILD_TAG-coverage'
                                sh 'docker rmi docker/cnab-app-base:$BUILD_TAG-coverage-experimental'
                            }
                            deleteDir()
                        }
                    }
                }
            }
        }
        stage('Test') {
            parallel {
                stage("Coverage") {
                    agent {
                        label 'ubuntu-1604-aufs-edge'
                    }
                    steps {
                        dir('src/github.com/docker/app') {
                            checkout scm
                            dir('_build') {
                                unstash "coverage-invocation-image"
                                sh 'docker load -i coverage-invocation-image.tar'
                            }
                            sh 'make -f docker.Makefile BUILD_TAG=$BUILD_TAG-coverage coverage'
                            archiveArtifacts '_build/ci-cov/all.out'
                            archiveArtifacts '_build/ci-cov/coverage.html'
                        }
                    }
                    post {
                        always {
                            sh 'docker rmi docker/cnab-app-base:$BUILD_TAG-coverage'
                            deleteDir()
                        }
                    }
                }
                stage("Coverage (experimental)") {
                    agent {
                        label 'ubuntu-1604-aufs-edge'
                    }
                    steps {
                        dir('src/github.com/docker/app') {
                            checkout scm
                            dir('_build') {
                                unstash "coverage-experimental-invocation-image"
                                sh 'docker load -i coverage-experimental-invocation-image.tar'
                            }
                            sh 'make EXPERIMENTAL=on -f docker.Makefile BUILD_TAG=$BUILD_TAG-coverage-experimental coverage'
                        }
                    }
                    post {
                        always {
                            sh 'docker rmi docker/cnab-app-base:$BUILD_TAG-coverage-experimental'
                            deleteDir()
                        }
                    }
                }
                stage("Gradle test") {
                    agent {
                        label 'ubuntu-1604-aufs-edge'
                    }
                    steps {
                        dir('src/github.com/docker/app') {
                            checkout scm
                            dir("bin") {
                                unstash "binaries"
                            }
                            sh 'make -f docker.Makefile gradle-test'
                        }
                    }
                    post {
                        always {
                            deleteDir()
                        }
                    }
                }
                stage("Test Linux") {
                    agent {
                        label 'ubuntu-1604-aufs-edge'
                    }
                    environment {
                        DOCKERAPP_BINARY = '../docker-app-linux'
                        DOCKERCLI_BINARY = '../docker-linux'
                    }
                    steps  {
                        dir('src/github.com/docker/app') {
                            checkout scm
                            dir('_build') {
                                unstash "invocation-image"
                                sh 'docker load -i invocation-image.tar'
                            }
                            unstash "binaries"
                            dir('examples') {
                                unstash "examples"
                            }
                            dir('e2e'){
                                unstash "e2e"
                            }
                            sh './docker-app-e2e-linux --e2e-path=e2e'
                        }
                    }
                    post {
                        always {
                            sh 'docker rmi docker/cnab-app-base:$BUILD_TAG'
                            deleteDir()
                        }
                    }
                }
            }
        }
    }
}
