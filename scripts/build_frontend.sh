curl -fsSL https://deb.nodesource.com/setup_17.x | sudo -E bash -

sudo apt-get update
sudo apt-get install -y apache2 build-essential nodejs

echo "Cloning frontend"
git clone https://github.com/Devops-2022-Group-R/itu-minitwit-frontend $HOME/frontend
cd $HOME/frontend

echo "VITE_API_URL=http://api.rhododevdron.swuwu.dk" >> .env

echo "Building frontend"
npm install
npm run build

echo "Allow apache"
sudo ufw allow 'Apache'

echo "Copy frontend to apache"
sudo cp -r $HOME/frontend/dist/* /var/www/html
