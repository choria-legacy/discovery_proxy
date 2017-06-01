desc "Set versions and create docs for a release"
task :prep_version do
  abort("Please specify CHORIA_VERSION") unless ENV["CHORIA_VERSION"]

  sh 'sed -i -re \'s/(.+"version": ").+/\1%s",/\' module/metadata.json' % ENV["CHORIA_VERSION"]
  sh 'sed -i -re \'s/(const version = ").+/\1%s"/\' cmd/cmd.go' % ENV["CHORIA_VERSION"]
  sh 'sed -i -re \'s/(discovery_proxy\/discovery_proxy-).+/\1%s",/\' module/manifests/init.pp' % ENV["CHORIA_VERSION"]

  changelog = File.readlines("CHANGELOG.md")

  File.open("CHANGELOG.md", "w") do |cl|
    changelog.each do |line|
      # rubocop:disable Metrics/LineLength
      cl.puts line

      if line =~ /^\|----------/
        cl.puts "|%s|      |Release %s                                                                                           |" % [Time.now.strftime("%Y/%m/%d"), ENV["CHORIA_VERSION"]]
      end
      # rubocop:enable Metrics/LineLength
    end
  end

  sh 'env GOOS=linux GOARCH=amd64 go build -o module/files/discovery_proxy-%s' % ENV["CHORIA_VERSION"]
  sh "git add CHANGELOG.md module cmd"
  sh "git commit -e -m '(misc) Release %s'" % ENV["CHORIA_VERSION"]
  sh "git tag %s" % ENV["CHORIA_VERSION"]
end

desc "Prepare and build the Puppet module"
task :release do
  rm Dir.glob("module/files/discovery_proxy*")
  Rake::Task[:prep_version].execute if ENV["CHORIA_VERSION"]

  Dir.chdir("module") do
    sh("/opt/puppetlabs/bin/puppet module build")
  end

  rm Dir.glob("module/files/discovery_proxy*")
end
