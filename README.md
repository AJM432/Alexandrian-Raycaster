![Alexandrian-Raycaster](https://github.com/user-attachments/assets/e94e726d-4f11-4d76-8343-c78f15ed4a3f)

# Alexandrian Raycaster
Simulates the Great Library of Alexandria with raycasting

[![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)](#)

## Demo
Coming soon

## Features
### Raycasting Methodology 
The program uses the perpendicular distance to the camera plane to prevent the fish-eye effect caused by using Euclidian distance. Rays are sent out from the player's direction in a 66° field of view. The resolution of the number of rays is the same as the number of horizontal pixels on the screen. The program contains a mini-map displaying the rays emanating from the player.
![image](https://github.com/user-attachments/assets/4e2df10f-1a9c-4fd8-8514-50701fc0ae16)

In green, the current player's position is shown; in yellow, the rays sent from the player are shown; in white, walls are represented; and in black, empty space is represented.

### Reading Ancient Books
The main focus of the program is the ability to read ancient books. The top 100 classical antiquity books (https://www.gutenberg.org/ebooks/bookshelf/24) were downloaded from a Project Gutenberg mirror site to allow for bulk downloading. The script used for the download is located in the `/get_books` subfolder. When the player gets within 2 block units of a wall, they are presented with a random book on the right side of the screen.
<p align="left">
<img src='https://github.com/user-attachments/assets/c16624dc-869c-4f83-97b4-106977e3f889' width='300'>
</p>

The program saves the user's current location in the book and subsequent visits to that book will remember the current location. Only one book is loaded at a time to prevent high memory usage.

## Feedback

If you have any feedback, please reach out to me at alvinjosematthew@gmail.com
## Acknowledgements
- [Library of Alexandria Wikipedia](https://en.wikipedia.org/wiki/Library_of_Alexandria)
- [Raycaster Guide](https://lodev.org/cgtutor/raycasting.html)


## License

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
## 🔗 Links
[![portfolio](https://img.shields.io/badge/my_portfolio-000?style=for-the-badge&logo=ko-fi&logoColor=white)](https://alvinmatthew.com/)
