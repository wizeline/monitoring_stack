
<!-- ABOUT THE PROJECT -->
## About The Project

This is a monitoring stack based on graphana, prometheus and alertmanager.

### Built With

This section should list any major frameworks/libraries used to bootstrap your project. Leave any add-ons/plugins for the acknowledgements section. Here are a few examples.

* [Docker](https://www.docker.com/)
* [Go Programming Language](https://go.dev)
* [Prometheus](https://prometheus.io/)
* [Alertmanager](https://prometheus.io/docs/alerting/latest/alertmanager/)
* [Grafana](https://grafana.com/)
* [MariaDB](https://mariadb.org/)



### Prerequisites

This code has been tested with the next versions of docker:
* docker (20.10.7) 
* docker compose (1.29.2)

It is expected to work with same or non-breaking-changes greater versions of the issued tools 

### Installation

1. Get a free API Key in order to use the openweathermap exporter [on openweathermap.org](https://home.openweathermap.org/users/sign_up)
2. Clone the repo
   ```sh
   git clone https://github.com/wizeline/monitoring_stack.git
   ```
3. copy environment template to .env file 
   ```sh
   cp .env_template .env
   ```
4. Enter your API and custom configuration of ports, user and passwords on this file `.env`
   
5. Run docker compose up --build


## Usage

- You can query on grafana on port 3000 (if you do not modify defaults) and import a dashboard from path /grafana/dashboards/weather.json to visualize real timem scraped data of temperature,humidity,pressure... 
The default user will be admin / admin 
- You can access the other resources on the specified ports on docker-compose.yml
<p align="right">(<a href="#top">back to top</a>)</p>

<!-- ROADMAP -->
## Roadmap

- [ ] Add Changelog
- [ ] Modify default login on Grafana
- [ ] Improve Grafana integration to autoload dashboards 
- [ ] Improve example preloaded data on prometheus 
- [ ] Whatever fancy feature you can think about (preload alerting, add new alerting methods, improve ux theme) 

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- LICENSE  TO REVIEW
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>
-->


<!-- OTHER INFO -->
## Other info & tips
If you want to add a new exporter to prometheus successfully do not forget to add its configuration also on `prometheus.yml` (scrape_configs:)

<!-- CONTACT -->
## Contact

Elena García - elena.garcia@wizeline.com

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* This is a fork from [ablx/monitoring_stack](https://github.com/ablx/monitoring_stack)

* The exporters for prometheus where also inspired on [dsmith73/openweathermap_exporter](https://github.com/dsmith73/openweathermap_exporter)

<p align="right">(<a href="#top">back to top</a>)</p>
