require 'digest'
require 'aws-sdk'

BUCKET_NAME = 'dickeyxxx_dev'

TARGETS = [
  {os: 'darwin', arch: 'amd64'},
  {os: 'linux', arch: 'amd64'},
  {os: 'linux', arch: '386'},
  {os: 'windows', arch: 'amd64'},
  {os: 'windows', arch: '386'}
]

VERSION = `./version.sh`.chomp

puts "hk VERSION: #{VERSION}"

task :build do
  FileUtils.mkdir_p 'dist'
  TARGETS.each do |target|
    path = "./dist/hk_#{target[:os]}_#{target[:arch]}"
    puts "building #{path}..."
    build(target[:os], target[:arch], path)
  end
end

task :gzip => :build do
  TARGETS.each do |target|
    path = "./dist/hk_#{target[:os]}_#{target[:arch]}"
    puts "gzipping #{path}..."
    system("gzip --keep -f #{path}")
    write_digest("#{path}.gz")
  end
end

namespace :deploy do
  task :release => :gzip do
    deploy(:release)
  end

  task :dev => :gzip do
    deploy(:dev)
  end
end

def deploy(channel)
  puts "deploying #{VERSION} to #{channel}..."
  bucket = get_s3_bucket
  TARGETS.each do |target|
    filename = "hk_#{target[:os]}_#{target[:arch]}.gz"
    local_path = "./dist/#{filename}"
    remote_path = "hk/#{channel}/#{VERSION}/#{filename}"
    remote_url = "#{BUCKET_NAME}.s3.amazonaws.com/#{remote_path}"
    puts "uploading #{local_path} to #{remote_url}"
    upload_file(bucket, local_path, remote_path)
    upload_file(bucket, local_path + ".sha1", remote_path + ".sha1")
  end
  version_path = "hk/#{channel}/VERSION"
  puts "setting to #{version_path} #{VERSION}"
  upload_file(bucket, 'VERSION', version_path)
end

def build(os, arch, path)
  ldflags = "-X main.VERSION #{VERSION}"
  args = "-o #{path} -ldflags \"#{ldflags}\""
  system("GOOS=#{os} GOARCH=#{arch} go build #{args}")
end

def write_digest(path)
  digest = Digest::SHA1.file(path).hexdigest
  File.open(path + '.sha1', 'w') { |f| f.write(digest) }
end

def get_s3_bucket
  s3 = AWS::S3.new
  s3.buckets[BUCKET_NAME]
end

def upload_file(bucket, local, remote)
  obj = bucket.objects[remote]
  obj.write(Pathname.new(local))
  obj.acl = :public_read
end

def upload_string(bucket, s, remote)
  obj = bucket.objects[remote]
  obj.write(s)
  obj.acl = :public_read
end
