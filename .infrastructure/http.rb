require 'net/http'
require 'json'

def droplet(id, token)
    uri = URI("https://api.digitalocean.com/v2/droplets/#{id}")
    req = Net::HTTP::Get.new(uri)
    req["Authorization"] = "Bearer #{token}"

    res = Net::HTTP.start(uri.hostname, uri.port, use_ssl: uri.scheme == 'https') { |http|
        http.request(req)
    }

    return JSON.parse(res.body)['droplet']
end