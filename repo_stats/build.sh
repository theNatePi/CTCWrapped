#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi
package_split=(${package//\// })
package_name=${package_split[$((${#package_split[@]}-1))]}
mkdir -p ./dist

platforms=(
  "darwin/amd64"
  "darwin/arm64"
  "linux/386"
  "linux/amd64"
  "linux/arm64"
  "windows/386"
  "windows/amd64"
  "windows/arm64"
  "linux/arm"
)

for platform in "${platforms[@]}"
do
  echo "Building for $platform"
  platform_split=(${platform//\// })
  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}
  output_name='./dist/'$package_name'-'$GOOS'-'$GOARCH
  if [ "$GOOS" = "windows" ]; then
    output_name+='.exe'
  fi

  env GOOS=$GOOS GOARCH=$GOARCH go build -o "$output_name" "$package"
  if [ $? -ne 0 ]; then
    echo 'An error has occurred! Aborting the script execution...'
    exit 1
  fi

  # Ensure downloadable artifacts have the expected permissions.
  if [ "$GOOS" = "windows" ]; then
    chmod 644 "$output_name"
  else
    chmod 755 "$output_name"
  fi
done
