# Coding Challenge
Imagine it’s 2049 and you are working for a company called SpaceTrouble that sends people to different places in our 
solar system. You are not the only one working in this industry. Your biggest competitor is a less-known company called 
SpaceX. Unfortunately, you both share the same launchpads and you cannot launch your rockets from the same place on the 
same day. There is a list of available launchpads and your spaceships go to places like Mars, Moon, Pluto, Asteroid 
Belt, Europa, Titan, and Ganymede. Every day you change the destination for all the launchpads. Every day of the week 
from the same launchpad has to be a “flight” to a different place.

Information about available launchpads and upcoming SpaceX launches you can find by SpaceX 
API: https://github.com/r-spacex/SpaceX-API

Your task is to create an API that will let your consumers book tickets online.

In order to do that you have to create 2 endpoints:

1. Endpoint to book a ticket where a client sends data like:
- First Name
- Last Name
- Gender
- Birthday
- Launchpad ID
- Destination ID
- Launch Date

You have to verify if the requested trip is possible on the day from the provided launchpad ID and does not overlap with 
SpaceX launches, if that’s the case then your flight is cancelled.

2. Endpoint to get all created Bookings.

Extra points:
- When you use docker/docker-compose to run the project.
- When you write unit/functional tests.
- When you create an endpoint to delete a booking.

Technical requirements:
- Please, use Golang and Postgres.
- Please, use GitHub or Bitbucket.
- Commit your changes often. Do not push the whole project in one commit.
