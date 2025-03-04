//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2021 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

// clusterintegrationtest acts as a test package to provide a compoenent test
// spanning multiple parts of the application, including everything that's
// required for a distributed setup. It thus acts like a mini "main" page and
// must be separated from the rest of the package to avoid circular import
// issues, etc.
package clusterintegrationtest
