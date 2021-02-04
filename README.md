<h1 align="center">The Kubectl Java Plugin</h1>

![version][go-shield]
![commit][commit-shield]
[![LICENSE][license-shield]][anti-996-url]

<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgements">Acknowledgements</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->

## About The Project

The kubectl java plugin aims to make you life easier while develop java application with k8s:

* find pods that running java application (👨🏻‍💻 working now)
* export jvm debug port and do port-forward directly
* more future...

### Built With

* [cobra](https://github.com/spf13/cobra)
* [color](https://github.com/fatih/color)
* [uitable](https://github.com/gosuri/uitable)

<!-- GETTING STARTED -->

## Getting Started

### Prerequisites

* Kubernetes Environment
    * [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl)
* Go sdk (optional)

### Installation

* For developer

```shell
git clone git@github.com:Ubisoft-potato/kubectl-java.git 
cd kubectl-java
go build
```

<!-- USAGE EXAMPLES -->

## Usage

* find pods that running java application in your cluster

```shell
kubectl-java list
```

**output will look like:**

```
context:dev	namespace:dev	maserURL:https://192.168.123.123:6443
NAME                                   	NODE      	STATUS 	CONTAINERS              	JDK
user-service-64d4f59c54-w9rwr          	dev-01	        Running	[user-service]          	openjdk version "1.8.0_232"
order-service-5654856bf6-9qb26     	dev-01	        Running	[order-service]     	        openjdk version "1.8.0_232"
chat-service-58fd5b4bf-gq8wz            dev-01	        Running	[chat-service]           	openjdk version "1.8.0_232"
```

<!-- ROADMAP -->

## Roadmap

See the [open issues](https://github.com/Ubisoft-potato/kubectl-java/issues) for a list of proposed features (and known
issues).



<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any
contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- LICENSE -->

## License

Distributed under the Anti 996-License-1.0 License. See `LICENSE` for more information.



<!-- CONTACT -->

## Contact

Project Link: [https://github.com/Ubisoft-potato/kubectl-java](https://github.com/Ubisoft-potato/kubectl-java)



<!-- ACKNOWLEDGEMENTS -->

## Acknowledgements

* [client-go](https://github.com/kubernetes/client-go)
* [cli-runtime](https://github.com/kubernetes/cli-runtime)

[go-shield]: https://img.shields.io/github/go-mod/go-version/Ubisoft-potato/kubectl-java

[commit-shield]: https://img.shields.io/github/last-commit/Ubisoft-potato/kubectl-java

[license-shield]: https://img.shields.io/badge/license-Anti%20996-blue.svg

[anti-996-url]: https://github.com/kattgu7/Anti-996-License